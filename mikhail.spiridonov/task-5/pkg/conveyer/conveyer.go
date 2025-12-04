package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	undefined = "undefined"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrChanFull     = errors.New("channel is full")
	ErrNoData       = errors.New("no data available")
)

type DefaultConveyer struct {
	size     int
	channels map[string]chan string
	handlers []handler
	mu       sync.RWMutex
	closed   bool
	running  bool
}

type handler interface {
	run(ctx context.Context) error
}

func New(size int) *DefaultConveyer {
	return &DefaultConveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]handler, 0),
		mu:       sync.RWMutex{},
		closed:   false,
		running:  false,
	}
}

func (c *DefaultConveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	c.running = true
	handlers := c.handlers
	c.mu.Unlock()

	defer c.close()

	group, groupCtx := errgroup.WithContext(ctx)

	for _, h := range handlers {
		handler := h

		group.Go(func() error {
			return handler.run(groupCtx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer stopped: %w", err)
	}

	return nil
}

func (c *DefaultConveyer) Send(input string, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channel, exists := c.channels[input]
	if !exists {
		return fmt.Errorf("%w: %s", ErrChanNotFound, input)
	}

	select {
	case channel <- data:
		return nil
	default:
		return ErrChanFull
	}
}

func (c *DefaultConveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channel, exists := c.channels[output]

	if !exists {
		return "", fmt.Errorf("%w: %s", ErrChanNotFound, output)
	}

	select {
	case data, ok := <-channel:
		if !ok {
			return undefined, nil
		}

		return data, nil
	default:
		return "", ErrNoData
	}
}

func (c *DefaultConveyer) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running || c.closed {
		return
	}

	inCh := c.getOrCreateChannelUnsafe(input)
	outCh := c.getOrCreateChannelUnsafe(output)

	c.handlers = append(c.handlers, &decorator{
		function: function,
		input:    inCh,
		output:   outCh,
	})
}

func (c *DefaultConveyer) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.running || c.closed {
		return
	}

	inChs := make([]chan string, len(inputs))
	for i, inputName := range inputs {
		inChs[i] = c.getOrCreateChannelUnsafe(inputName)
	}

	outCh := c.getOrCreateChannelUnsafe(output)

	c.handlers = append(c.handlers, &multiplexer{
		function: function,
		input:    inChs,
		output:   outCh,
	})
}

func (c *DefaultConveyer) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.running || c.closed {
		return
	}

	inCh := c.getOrCreateChannelUnsafe(input)

	outChs := make([]chan string, len(outputs))
	for i, outputName := range outputs {
		outChs[i] = c.getOrCreateChannelUnsafe(outputName)
	}

	c.handlers = append(c.handlers, &separator{
		function: function,
		input:    inCh,
		output:   outChs,
	})
}

func (c *DefaultConveyer) getOrCreateChannelUnsafe(name string) chan string {
	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *DefaultConveyer) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.closed = true
	for _, channel := range c.channels {
		close(channel)
	}
}
