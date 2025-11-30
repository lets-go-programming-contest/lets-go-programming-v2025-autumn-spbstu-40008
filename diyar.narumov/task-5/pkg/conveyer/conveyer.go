package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChannelNotFound = errors.New("chan not found")

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	chMutex  sync.RWMutex
	handlers []func(ctx context.Context) error
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		chMutex:  sync.RWMutex{},
		handlers: []func(context.Context) error{},
	}
}

func (c *conveyerImpl) ensureChannel(name string) chan string {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *conveyerImpl) getChannel(name string) chan string {
	c.chMutex.RLock()
	defer c.chMutex.RUnlock()

	return c.channels[name]
}

func (c *conveyerImpl) RegisterDecorator(
	handlerFn func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inputCh := c.ensureChannel(input)
	outputCh := c.ensureChannel(output)

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputCh, outputCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	handlerFn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputsChs := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputsChs[i] = c.ensureChannel(name)
	}

	outputCh := c.ensureChannel(output)

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputsChs, outputCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	handlerFn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputCh := c.ensureChannel(input)
	outputsChs := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputsChs[i] = c.ensureChannel(name)
	}

	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputCh, outputsChs)
	})
}

func (c *conveyerImpl) Send(input, data string) error {
	inputCh := c.getChannel(input)
	if inputCh == nil {
		return ErrChannelNotFound
	}

	inputCh <- data

	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	outputCh := c.getChannel(output)
	if outputCh == nil {
		return "", ErrChannelNotFound
	}

	val, ok := <-outputCh
	if !ok {
		return "undefined", nil
	}

	return val, nil
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)

	var waitGroup sync.WaitGroup
	for _, handler := range c.handlers {
		waitGroup.Add(1)

		go func(h func(ctx context.Context) error) {
			defer waitGroup.Done()

			if err := h(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(handler)
	}

	var err error
	select {
	case err = <-errCh:
		cancel()
	case <-ctx.Done():
	}

	waitGroup.Wait()
	c.closeAllChannels()

	return err
}

func (c *conveyerImpl) closeAllChannels() {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}
