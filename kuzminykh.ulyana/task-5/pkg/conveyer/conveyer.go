package conveyer

import (
	"context"
	"errors"
	"fmt"
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
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exist := c.channels[name]; exist {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel
	return channel
}

func (c *conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	channel, exist := c.channels[name]

	return channel, exist
}

func (c *conveyer) RegisterDecorator(
	handlerFunc func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inCh := c.createChannel(input)
	outCh := c.createChannel(output)

	handler := func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outCh)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()
}

func (c *conveyer) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inCh := make([]chan string, len(inputs))
	for i, name := range inputs {
		inCh[i] = c.createChannel(name)
	}
	outCh := c.createChannel(output)

	handler := func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outCh)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()
}

func (c *conveyer) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.createChannel(input)
	outCh := make([]chan string, len(outputs))
	for i, name := range outputs {
		outCh[i] = c.createChannel(name)
	}
	handler := func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outCh)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()
}

func (c *conveyer) Send(input string, data string) error {
	channel, exist := c.getChannel(input)
	if !exist {
		return ErrNotFound
	}

	channel <- data

	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	channel, exist := c.getChannel(output)
	if !exist {
		return "", ErrNotFound
	}

	val, exist := <-channel
	if !exist {
		return "undefined", nil
	}

	return val, nil
}

func (c *conveyer) Run(ctx context.Context) error {
	c.mu.RLock()
	handlers := make([]func(context.Context) error, len(c.handlers))
	copy(handlers, c.handlers)
	c.mu.RUnlock()

	var waitGroup sync.WaitGroup
	errCh := make(chan error, 1)

	for _, handler := range handlers {
		waitGroup.Add(1)
		go func(handlerFunc func(context.Context) error) {
			defer waitGroup.Done()
			if err := handlerFunc(ctx); err != nil {
				select {
				case errCh <- fmt.Errorf("handler error: %w", err):
				default:
				}
			}
		}(handler)
	}

	done := make(chan struct{})
	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		<-done
		c.closeChannels()
		return err
	case <-done:
		c.closeChannels()
		return nil
	case <-ctx.Done():
		<-done
		c.closeChannels()
		return fmt.Errorf("context canceled: %w", ctx.Err())
	}
}

func (c *conveyer) closeChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, channel := range c.channels {
		close(channel)
	}
}
