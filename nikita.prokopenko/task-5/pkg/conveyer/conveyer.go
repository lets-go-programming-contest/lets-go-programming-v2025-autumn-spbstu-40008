// pkg/conveyer/conveyer.go
package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mu       sync.RWMutex
	streams  map[string]chan string
	actions  []func(ctx context.Context) error
	capacity int
}

func New(bufferSize int) *Conveyer {
	return &Conveyer{
		streams:  make(map[string]chan string),
		actions:  make([]func(ctx context.Context) error, 0),
		capacity: bufferSize,
	}
}

func (c *Conveyer) acquireChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.streams[name]; exists {
		return ch
	}

	ch := make(chan string, c.capacity)
	c.streams[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	processor func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	input := c.acquireChannel(inputName)
	output := c.acquireChannel(outputName)

	c.mu.Lock()
	c.actions = append(c.actions, func(ctx context.Context) error {
		return processor(ctx, input, output)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	processor func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	inputs := make([]chan string, len(inputNames))
	for i, name := range inputNames {
		inputs[i] = c.acquireChannel(name)
	}
	output := c.acquireChannel(outputName)

	c.mu.Lock()
	c.actions = append(c.actions, func(ctx context.Context) error {
		return processor(ctx, inputs, output)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	processor func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	input := c.acquireChannel(inputName)
	outputs := make([]chan string, len(outputNames))
	for i, name := range outputNames {
		outputs[i] = c.acquireChannel(name)
	}

	c.mu.Lock()
	c.actions = append(c.actions, func(ctx context.Context) error {
		return processor(ctx, input, outputs)
	})
	c.mu.Unlock()
}

func (c *Conveyer) Send(channelName string, data string) error {
	c.mu.RLock()
	ch, exists := c.streams[channelName]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	c.mu.RLock()
	ch, exists := c.streams[channelName]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	value, isOpen := <-ch
	if !isOpen {
		return "undefined", nil
	}

	return value, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	for _, action := range c.actions {
		wg.Add(1)
		go func(act func(ctx context.Context) error) {
			defer wg.Done()
			if err := act(ctx); err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(action)
	}

	wg.Wait()

	c.mu.Lock()
	for _, ch := range c.streams {
		close(ch)
	}
	c.mu.Unlock()

	return firstErr
}