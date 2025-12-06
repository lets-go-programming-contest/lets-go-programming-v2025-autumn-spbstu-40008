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
	chMutex  *sync.RWMutex
	handlers []func(ctx context.Context) error
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		chMutex:  &sync.RWMutex{},
		handlers: make([]func(context.Context) error, 0),
	}
}

func (c *conveyerImpl) ensureChannel(name string) chan string {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *conveyerImpl) getChannel(name string) chan string {
	c.chMutex.RLock()
	defer c.chMutex.RUnlock()

	return c.channels[name]
}

func (c *conveyerImpl) RegisterDecorator(
	handlerFunc func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inCh := c.ensureChannel(input)
	outCh := c.ensureChannel(output)

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputChannels[i] = c.ensureChannel(name)
	}

	outCh := c.ensureChannel(output)

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannels, outCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.ensureChannel(input)
	outputChannels := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputChannels[i] = c.ensureChannel(name)
	}

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inCh, outputChannels)
	})
}

func (c *conveyerImpl) Send(input, data string) error {
	channel := c.getChannel(input)
	if channel == nil {
		return ErrChannelNotFound
	}

	channel <- data

	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	channel := c.getChannel(output)
	if channel == nil {
		return "", ErrChannelNotFound
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		currentHandler := handler

		errGroup.Go(func() error {
			return currentHandler(ctx)
		})
	}

	err := errGroup.Wait()

	c.closeAllChannels()

	return err
}

func (c *conveyerImpl) closeAllChannels() {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}
