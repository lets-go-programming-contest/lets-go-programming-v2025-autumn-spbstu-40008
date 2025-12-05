package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mutex    sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		mutex:    sync.RWMutex{},
		handlers: []func(ctx context.Context) error{},
	}
}

func (c *Conveyer) register(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) RegisterDecorator(
	callback func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inCh := c.register(input)
	outCh := c.register(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return callback(ctx, inCh, outCh)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	callback func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = c.register(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return callback(ctx, inChans, c.register(output))
	})
}

func (c *Conveyer) RegisterSeparator(
	callback func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = c.register(name)
	}

	inCh := c.register(input)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return callback(ctx, inCh, outChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer func() {
		c.mutex.RLock()
		defer c.mutex.RUnlock()

		for _, ch := range c.channels {
			close(ch)
		}
	}()

	errgr, ctx := errgroup.WithContext(ctx)

	c.mutex.RLock()

	for _, handler := range c.handlers {
		errgr.Go(func() error {
			return handler(ctx)
		})
	}

	c.mutex.RUnlock()

	if err := errgr.Wait(); err != nil {
		return fmt.Errorf("run pipeline: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.mutex.RLock()

	channel, ok := c.channels[input]

	c.mutex.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mutex.RLock()

	channel, ok := c.channels[output] //nolint:varnamelen

	c.mutex.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return data, nil
}
