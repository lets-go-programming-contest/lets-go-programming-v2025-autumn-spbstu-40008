package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("chan not found")

type Conveyer interface {
	RegisterDecorator(
		handler func(ctx context.Context, input, output chan string) error,
		input, output string,
	)
	RegisterMultiplexer(
		handler func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		handler func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	chMutex  sync.RWMutex
	handlers []func(ctx context.Context) error
}

func New(size int) Conveyer {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
		chMutex:  sync.RWMutex{},
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
	handler func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inputChan := c.ensureChannel(input)
	outputChan := c.ensureChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChan, outputChan)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputChans[i] = c.ensureChannel(name)
	}
	outputChan := c.ensureChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChans, outputChan)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := c.ensureChannel(input)
	outputChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outputChans[i] = c.ensureChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChan, outputChans)
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
	group, ctx := errgroup.WithContext(ctx)

	for _, handlerFunc := range c.handlers {
		group.Go(func() error { return handlerFunc(ctx) })
	}

	err := group.Wait()
	c.closeAllChannels()

	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

func (c *conveyerImpl) closeAllChannels() {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}