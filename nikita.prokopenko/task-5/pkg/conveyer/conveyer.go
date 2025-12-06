package conveyer

import (
	"context"
	"errors"
	"fmt"
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
		handlers: make([]func(ctx context.Context) error, 0),
		buffer:   size,
	}
}

func (c *Conveyer) createOrGetChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if chnl, ok := c.channels[name]; ok {
		return chnl
	}

	chnl := make(chan string, c.buffer)
	c.channels[name] = chnl
	return chnl
}

func (c *Conveyer) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	c.createOrGetChannel(inputName)
	c.createOrGetChannel(outputName)

	handlerFn := func(ctx context.Context) error {
		input := c.createOrGetChannel(inputName)
		output := c.createOrGetChannel(outputName)
		return handler(ctx, input, output)
	}

	c.handlers = append(c.handlers, handlerFn)
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

	handlerFn := func(ctx context.Context) error {
		inputs := make([]chan string, len(inputNames))
		for i, name := range inputNames {
			inputs[i] = c.createOrGetChannel(name)
		}

		output := c.createOrGetChannel(outputName)
		return handler(ctx, inputs, output)
	}

	c.handlers = append(c.handlers, handlerFn)
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

	handlerFn := func(ctx context.Context) error {
		input := c.createOrGetChannel(inputName)
		outputs := make([]chan string, len(outputNames))
		for i, name := range outputNames {
			outputs[i] = c.createOrGetChannel(name)
		}

		return handler(ctx, input, outputs)
	}

	c.handlers = append(c.handlers, handlerFn)
}

func (c *Conveyer) Send(channelName string, data string) error {
	c.mu.RLock()
	channel, exists := c.channels[channelName]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("%w: channel %s not found", ErrChanNotFound, channelName)
	}

	select {
	case channel <- data:
		return nil
	default:
		channel <- data
		return nil
	}
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[channelName]
	c.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("%w: channel %s not found", ErrChanNotFound, channelName)
	}

	value, ok := <-channel
	if !ok {
		return "", nil
	}

	return value, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errChan := make(chan error, 1)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, handler := range c.handlers {
		waitGroup.Add(1)

		go func(h func(context.Context) error) {
			defer waitGroup.Done()
			if err := h(runCtx); err != nil {
				select {
				case errChan <- fmt.Errorf("handler error: %w", err):
				default:
				}
				cancel()
			}
		}(handler)
	}

	go func() {
		waitGroup.Wait()
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	case err := <-errChan:
		if err != nil {
			return err
		}
	}

	return nil
}