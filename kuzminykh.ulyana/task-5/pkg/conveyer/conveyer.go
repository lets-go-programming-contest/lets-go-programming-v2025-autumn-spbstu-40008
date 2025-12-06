package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("chan not found")

type conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(context.Context) error

	mu sync.RWMutex
}

func New(size int) *conveyer {
	return &conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
		mu:       sync.RWMutex{},
	}
}

func (c *conveyer) createChannel(name string) chan string {
	c.mu.RLock()
	if channel, ok := c.channels[name]; ok {
		c.mu.RUnlock()

		return channel
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if channel, ok := c.channels[name]; ok {

		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	channel, ok := c.channels[name]

	return channel, ok
}

func (c *conveyer) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inCh := c.createChannel(input)
	outCh := c.createChannel(output)

	handler := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()
}

func (c *conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	inCh := make([]chan string, len(inputs))
	for i, name := range inputs {
		inCh[i] = c.createChannel(name)
	}
	outCh := c.createChannel(output)

	handler := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()

}

func (c *conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	inCh := c.createChannel(input)
	outCh := make([]chan string, len(outputs))
	for i, name := range outputs {
		outCh[i] = c.createChannel(name)
	}
	handler := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()
}

func (c *conveyer) Send(input string, data string) error {
	channel, ok := c.getChannel(input)
	if !ok {

		return ErrNotFound
	}

	channel <- data

	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	channel, ok := c.getChannel(output)
	if !ok {
		return "", ErrNotFound
	}

	val, ok := <-channel
	if !ok {

		return "undefined", nil
	}

	return val, nil
}

func (c *conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c.mu.RLock()
	handlers := make([]func(context.Context) error, len(c.handlers))
	copy(handlers, c.handlers)
	c.mu.RUnlock()

	var waitGroup sync.WaitGroup
	errCh := make(chan error, 1)

	for _, handler := range handlers {
		waitGroup.Add(1)
		go func(t func(context.Context) error) {
			defer waitGroup.Done()
			if err := t(ctx); err != nil {
				select {
				case errCh <- err:
					cancel()
				default:
				}
			}
		}(handler)
	}

	go func() {
		waitGroup.Wait()

		close(errCh)
	}()

	select {
	case err, ok := <-errCh:
		if ok && err != nil {

			return err
		}

		return nil
	case <-ctx.Done():

		return ctx.Err()
	}
}
