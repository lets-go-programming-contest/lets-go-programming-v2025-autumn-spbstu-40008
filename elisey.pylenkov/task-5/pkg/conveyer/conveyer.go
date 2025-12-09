package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(ctx context.Context) error, 0),
		mu:       sync.RWMutex{},
	}
}

func (c *Conveyer) getOrCreateChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) getChan(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[name]

	return ch, ok
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inChan := c.getOrCreateChan(input)
	outChan := c.getOrCreateChan(output)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inChan, outChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = c.getOrCreateChan(name)
	}

	outChan := c.getOrCreateChan(output)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inChans, outChan)
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inChan := c.getOrCreateChan(input)
	outChans := make([]chan string, len(outputs))

	for i, name := range outputs {
		outChans[i] = c.getOrCreateChan(name)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inChan, outChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		currentHandler := handler

		group.Go(func() error {
			return currentHandler(groupCtx)
		})
	}

	if err := group.Wait(); err != nil {
		c.closeAllChannels()
		return err
	}

	c.closeAllChannels()

	return nil
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	ch, exists := c.getChan(input)
	if !exists {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, exists := c.getChan(output)
	if !exists {
		return "", ErrChanNotFound
	}

	val, isOpen := <-ch
	if !isOpen {
		return "undefined", nil
	}

	return val, nil
}
