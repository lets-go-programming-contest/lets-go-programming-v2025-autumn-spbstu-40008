package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("chan not found")

const UndefinedValue = "undefined"

type handlerFunc func(ctx context.Context) error

type Conveyer struct {
	size     int
	mutex    *sync.RWMutex
	channels map[string]chan string
	handlers []handlerFunc
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		mutex:    &sync.RWMutex{},
		channels: make(map[string]chan string),
		handlers: make([]handlerFunc, 0),
	}
}

func (c *Conveyer) ensureChannel(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, <-chan string, chan<- string) error,
	input, output string,
) {
	in := c.ensureChannel(input)
	out := c.ensureChannel(output)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []<-chan string, chan<- string) error,
	inputs []string,
	output string,
) {
	in := make([]<-chan string, len(inputs))
	for i, name := range inputs {
		in[i] = c.ensureChannel(name)
	}
	out := c.ensureChannel(output)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, <-chan string, []chan<- string) error,
	input string,
	outputs []string,
) {
	in := c.ensureChannel(input)
	out := make([]chan<- string, len(outputs))
	for i, name := range outputs {
		out[i] = c.ensureChannel(name)
	}
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) Send(pipe, value string) error {
	c.mutex.RLock()
	channel, ok := c.channels[pipe]
	c.mutex.RUnlock()
	if !ok {
		return ErrChannelNotFound
	}
	channel <- value
	return nil
}

func (c *Conveyer) Recv(pipe string) (string, error) {
	c.mutex.RLock()
	channel, ok := c.channels[pipe]
	c.mutex.RUnlock()
	if !ok {
		return "", ErrChannelNotFound
	}
	value, open := <-channel
	if !open {
		return UndefinedValue, nil
	}
	return value, nil
}

func (c *Conveyer) closeAll() {
	c.mutex.Lock()
	for _, channel := range c.channels {
		close(channel)
	}
	c.mutex.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, h := range c.handlers {
		fn := h
		g.Go(func() error {
			return fn(ctx)
		})
	}
	err := g.Wait()
	c.closeAll()
	if err != nil {
		return fmt.Errorf("run failed: %w", err)
	}
	return nil
}

