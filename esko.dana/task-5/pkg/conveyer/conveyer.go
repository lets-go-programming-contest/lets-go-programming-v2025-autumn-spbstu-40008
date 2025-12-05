package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("chan not found")

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	chMutex  sync.RWMutex
	handlers []func(ctx context.Context) error
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
	}
}

func (c *conveyerImpl) ensureChannel(name string) chan string {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) getChannel(name string) chan string {
	c.chMutex.RLock()
	defer c.chMutex.RUnlock()
	return c.channels[name]
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inCh := c.ensureChannel(input)
	outCh := c.ensureChannel(output)

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inChs := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChs[i] = c.ensureChannel(name)
	}
	outCh := c.ensureChannel(output)

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.ensureChannel(input)
	outChs := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChs[i] = c.ensureChannel(name)
	}

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	})
}

func (c *conveyerImpl) Send(input, data string) error {
	ch := c.getChannel(input)
	if ch == nil {
		return ErrChannelNotFound
	}
	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch := c.getChannel(output)
	if ch == nil {
		return "", ErrChannelNotFound
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return data, nil
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		h := handler
		g.Go(func() error {
			return h(ctx)
		})
	}

	err := g.Wait()
	c.closeAllChannels()
	return err
}

func (c *conveyerImpl) closeAllChannels() {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()
	for _, ch := range c.channels {
		close(ch)
	}
}
