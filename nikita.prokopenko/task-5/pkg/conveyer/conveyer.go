// pkg/conveyer/conveyer.go
package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	buffer   int
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
		buffer:   size,
	}
}

func (c *Conveyer) createOrGetChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.buffer)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[name]
	return ch, ok
}

func (c *Conveyer) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	c.createOrGetChannel(inputName)
	c.createOrGetChannel(outputName)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		input := c.createOrGetChannel(inputName)
		output := c.createOrGetChannel(outputName)
		return handler(ctx, input, output)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	for _, name := range inputNames {
		c.createOrGetChannel(name)
	}
	c.createOrGetChannel(outputName)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		inputs := make([]chan string, len(inputNames))
		for i, name := range inputNames {
			inputs[i] = c.createOrGetChannel(name)
		}
		output := c.createOrGetChannel(outputName)
		return handler(ctx, inputs, output)
	})
}

func (c *Conveyer) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	c.createOrGetChannel(inputName)
	for _, name := range outputNames {
		c.createOrGetChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		input := c.createOrGetChannel(inputName)
		outputs := make([]chan string, len(outputNames))
		for i, name := range outputNames {
			outputs[i] = c.createOrGetChannel(name)
		}
		return handler(ctx, input, outputs)
	})
}

func (c *Conveyer) Send(channelName string, data string) error {
	ch, ok := c.getChannel(channelName)
	if !ok {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	ch, ok := c.getChannel(channelName)
	if !ok {
		return "", ErrChanNotFound
	}

	msg, isOpen := <-ch
	if !isOpen {
		return "undefined", nil
	}

	return msg, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	for _, handler := range c.handlers {
		wg.Add(1)
		go func(h func(context.Context) error) {
			defer wg.Done()
			if err := h(ctx); err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(handler)
	}

	wg.Wait()

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, ch := range c.channels {
		close(ch)
	}

	return firstErr
}