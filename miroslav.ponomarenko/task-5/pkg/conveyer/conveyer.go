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
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	in := c.getChannel(input)
	out := c.getChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, in, out)
	}

	c.workers = append(c.workers, task)
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var ins []chan string
	for _, name := range inputs {
		ins = append(ins, c.getChannel(name))
	}
	out := c.getChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, ins, out)
	}

	c.workers = append(c.workers, task)
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	in := c.getChannel(input)
	var outs []chan string
	for _, name := range outputs {
		outs = append(outs, c.getChannel(name))
	}

	task := func(ctx context.Context) error {
		return fn(ctx, in, outs)
	}

	c.workers = append(c.workers, task)
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.cleanup()

	g, gCtx := errgroup.WithContext(ctx)
	c.mutex.RLock()
	tasks := make([]func(context.Context) error, len(c.workers))
	copy(tasks, c.workers)
	c.mutex.RUnlock()

	for _, t := range tasks {
		task := t
		g.Go(func() error {
			return task(gCtx)
		})
	}

	if err := g.Wait(); err != nil {
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
	ch, ok := c.channels[input]
	c.mutex.RUnlock()

	if !ok {
		return ErrChannelNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mutex.RLock()
	ch, ok := c.channels[output]
	c.mutex.RUnlock()

	if !ok {
		return "", ErrChannelNotFound
	}

	val, ok := <-ch
	if !ok {
		return UndefinedValue, nil
	}

	return val, nil
}
