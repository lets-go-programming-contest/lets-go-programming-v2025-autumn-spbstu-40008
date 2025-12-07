package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const UndefinedValue = "undefined"

var ErrChannelNotFound = errors.New("chan not found")

type Conveyer struct {
	size     int
	channels map[string]chan string
	workers  []func(ctx context.Context) error
	mutex    sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]func(ctx context.Context) error, 0),
		mutex:    sync.RWMutex{},
	}
}

func (c *Conveyer) getChannel(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.getChannel(input)
	outChannel := c.getChannel(output)

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outChannel)
	}

	c.workers = append(c.workers, task)
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ins := make([]chan string, len(inputs))
	for i, name := range inputs {
		ins[i] = c.getChannel(name)
	}

	outChannel := c.getChannel(output)

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, ins, outChannel)
	}

	c.workers = append(c.workers, task)
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.getChannel(input)
	outs := make([]chan string, len(outputs))
	for i, name := range outputs {
		outs[i] = c.getChannel(name)
	}

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outs)
	}

	c.workers = append(c.workers, task)
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.cleanup()

	errGroup, gCtx := errgroup.WithContext(ctx)

	c.mutex.RLock()
	tasks := make([]func(context.Context) error, len(c.workers))
	copy(tasks, c.workers)
	c.mutex.RUnlock()

	for _, t := range tasks {
		task := t

		errGroup.Go(func() error {
			return task(gCtx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer error: %w", err)
	}

	return nil
}

func (c *Conveyer) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mutex.RLock()
	channel, exists := c.channels[input]
	c.mutex.RUnlock()

	if !exists {
		return ErrChannelNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mutex.RLock()
	channel, exists := c.channels[output]
	c.mutex.RUnlock()

	if !exists {
		return "", ErrChannelNotFound
	}

	val, exists := <-channel
	if !exists {
		return UndefinedValue, nil
	}

	return val, nil
}
