package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("chan not found")

// Conveyer defines the interface for the conveyor belt.
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

// New creates a new instance of Conveyer.
func New(size int) Conveyer {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
		chMutex:  sync.RWMutex{},
	}
}

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	chMutex  sync.RWMutex
	handlers []func(ctx context.Context) error
}

func (c *conveyerImpl) ensureChannel(name string) chan string {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	if existingChan, exists := c.channels[name]; exists {
		return existingChan
	}

	newChan := make(chan string, c.size)
	c.channels[name] = newChan

	return newChan
}

func (c *conveyerImpl) getChannel(name string) chan string {
	c.chMutex.RLock()
	defer c.chMutex.RUnlock()

	return c.channels[name]
}

func (c *conveyerImpl) RegisterDecorator(
	handler func(ctx context.Context, inputChan, outputChan chan string) error,
	input, output string,
) {
	inputChan := c.ensureChannel(input)
	outputChan := c.ensureChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChan, outputChan)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	handler func(ctx context.Context, inputChans []chan string, outputChan chan string) error,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputChans[i] = c.ensureChannel(name)
	}

	outputChan := c.ensureChannel(output) // Разделено пустой строкой

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChans, outputChan)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	handler func(ctx context.Context, inputChan chan string, outputChans []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := c.ensureChannel(input)
	outputChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outputChans[i] = c.ensureChannel(name)
	}

	outputChan := c.ensureChannel(output) // Разделено пустой строкой

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, inputChan, outputChans)
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
	group, innerCtx := errgroup.WithContext(ctx)

	for _, handlerFunc := range c.handlers {
		group.Go(func() error {
			return handlerFunc(innerCtx)
		})
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

