package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Conveyer struct {
	mu        sync.RWMutex
	streams   map[string]chan string
	processes []func(ctx context.Context) error
	bufSize   int
}

func New(size int) *Conveyer {
	return &Conveyer{
		streams:   make(map[string]chan string),
		processes: make([]func(ctx context.Context) error, 0),
		bufSize:   size,
	}
}

func (c *Conveyer) ensureChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.streams[name]; ok {
		return ch
	}
	ch := make(chan string, c.bufSize)
	c.streams[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	in, out string,
) {
	c.ensureChan(in)
	c.ensureChan(out)

	c.processes = append(c.processes, func(ctx context.Context) error {
		input := c.ensureChan(in)
		output := c.ensureChan(out)
		return fn(ctx, input, output)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	ins []string,
	out string,
) {
	for _, name := range ins {
		c.ensureChan(name)
	}
	c.ensureChan(out)

	c.processes = append(c.processes, func(ctx context.Context) error {
		inputs := make([]chan string, len(ins))
		for i, name := range ins {
			inputs[i] = c.ensureChan(name)
		}
		return fn(ctx, inputs, c.ensureChan(out))
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	in string,
	outs []string,
) {
	c.ensureChan(in)
	for _, name := range outs {
		c.ensureChan(name)
	}

	c.processes = append(c.processes, func(ctx context.Context) error {
		outputs := make([]chan string, len(outs))
		for i, name := range outs {
			outputs[i] = c.ensureChan(name)
		}
		return fn(ctx, c.ensureChan(in), outputs)
	})
}

func (c *Conveyer) Send(pipe string, data string) error {
	c.mu.RLock()
	ch, ok := c.streams[pipe]
	c.mu.RUnlock()

	if !ok {
		return errors.New("chan not found")
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(pipe string) (string, error) {
	c.mu.RLock()
	ch, ok := c.streams[pipe]
	c.mu.RUnlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, proc := range c.processes {
		procCopy := proc
		g.Go(func() error {
			return procCopy(ctx)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
