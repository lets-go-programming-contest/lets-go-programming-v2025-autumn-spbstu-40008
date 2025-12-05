package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("chan not found")
)

type conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func() error

	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(size int) *conveyer {
	return &conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func() error, 0),
	}
}

func (c *conveyer) createChannel(name string) chan string {
	c.mu.RLock()
	if ch, ok := c.channels[name]; ok {
		c.mu.RUnlock()
		return ch
	}
	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.channels[name]
	return ch, ok
}

func (c *conveyer) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error, input, output string) {
	c.handlers = append(c.handlers, func() error {
		in := c.createChannel(input)
		out := c.createChannel(output)
		return fn(c.ctx, in, out)
	})
}

func (c *conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	c.handlers = append(c.handlers, func() error {
		in := make([]chan string, len(inputs))
		for i, name := range inputs {
			in[i] = c.createChannel(name)
		}
		out := c.createChannel(output)
		return fn(c.ctx, in, out)
	})
}

func (c *conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	c.handlers = append(c.handlers, func() error {
		in := c.createChannel(input)
		out := make([]chan string, len(outputs))
		for i, name := range outputs {
			out[i] = c.createChannel(name)
		}
		return fn(c.ctx, in, out)
	})
}

func (c *conveyer) Send(input string, data string) error {
	ch, ok := c.getChannel(input)
	if !ok {
		return ErrNotFound
	}

	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case ch <- data:
		return nil
	}
}

func (c *conveyer) Recv(output string) (string, error) {
	ch, ok := c.getChannel(output)
	if !ok {
		return "", ErrNotFound
	}

	select {
	case <-c.ctx.Done():
		return "", c.ctx.Err()
	case val, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return val, nil
	}
}

func (c *conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()

	errCh := make(chan error, 1)

	for _, handler := range c.handlers {
		c.wg.Add(1)
		go func(t func() error) {
			defer c.wg.Done()
			if err := t(); err != nil {
				select {
				case errCh <- err:
					c.cancel()
				default:
				}
			}
		}(handler)
	}

	go func() {
		c.wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
}
