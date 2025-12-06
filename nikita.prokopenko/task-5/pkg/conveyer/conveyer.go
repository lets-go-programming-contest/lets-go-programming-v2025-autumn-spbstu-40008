// pkg/conveyer/conveyer.go
package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channels map[string]chan string
	actions  []func(ctx context.Context) error
	bufSize  int
	mu       sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		actions:  []func(ctx context.Context) error{},
		bufSize:  size,
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.bufSize)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, in, out chan string) error,
	inName, outName string,
) {
	inCh := c.getOrCreateChannel(inName)
	outCh := c.getOrCreateChannel(outName)

	c.mu.Lock()
	c.actions = append(c.actions, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, ins []chan string, out chan string) error,
	inNames []string,
	outName string,
) {
	inChs := make([]chan string, len(inNames))
	for i, name := range inNames {
		inChs[i] = c.getOrCreateChannel(name)
	}
	outCh := c.getOrCreateChannel(outName)

	c.mu.Lock()
	c.actions = append(c.actions, func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, in chan string, outs []chan string) error,
	inName string,
	outNames []string,
) {
	inCh := c.getOrCreateChannel(inName)
	outChs := make([]chan string, len(outNames))
	for i, name := range outNames {
		outChs[i] = c.getOrCreateChannel(name)
	}

	c.mu.Lock()
	c.actions = append(c.actions, func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	})
	c.mu.Unlock()
}

func (c *Conveyer) Send(name string, data string) error {
	c.mu.RLock()
	ch, exists := c.channels[name]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[name]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return val, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var errOnce sync.Once
	var firstErr error

	for _, action := range c.actions {
		wg.Add(1)
		go func(act func(ctx context.Context) error) {
			defer wg.Done()
			if err := act(ctx); err != nil {
				errOnce.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(action)
	}

	wg.Wait()

	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mu.Unlock()

	return firstErr
}