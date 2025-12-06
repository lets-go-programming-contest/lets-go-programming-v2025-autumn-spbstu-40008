package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mu        sync.RWMutex
	channels  map[string]chan string
	handlers  []func(ctx context.Context) error
	buffer    int
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:       sync.RWMutex{},
		channels: make(map[string]chan string),
		handlers: make([]func(ctx context.Context) error, 0),
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
	c.mu.RLock()
	ch, exists := c.channels[channelName]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[channelName]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	value, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return value, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	for _, handler := range c.handlers {
		wg.Add(1)
		go func(h func(ctx context.Context) error) {
			defer wg.Done()
			if err := h(ctx); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}(handler)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("handler execution failed: %w", err)
		}
	}

	return nil
}