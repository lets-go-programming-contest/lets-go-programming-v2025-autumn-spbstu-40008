package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mu        sync.RWMutex
	channels  map[string]chan string
	handlers  []func(ctx context.Context) error
	bufferLen int
}

func New(buffer int) *Conveyer {
	return &Conveyer{
		mu:        sync.RWMutex{},
		channels:  make(map[string]chan string),
		handlers:  make([]func(ctx context.Context) error, 0),
		bufferLen: buffer,
	}
}

func (c *Conveyer) ensureChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channelRef, ok := c.channels[name]; ok {
		return channelRef
	}

	channelRef := make(chan string, c.bufferLen)
	c.channels[name] = channelRef

	return channelRef
}

func (c *Conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channelRef, ok := c.channels[name]

	return channelRef, ok
}

func (c *Conveyer) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	inputChannel := c.ensureChannel(inputName)
	outputChannel := c.ensureChannel(outputName)

	task := func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannel)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	inputChannels := make([]chan string, 0, len(inputNames))
	for _, name := range inputNames {
		inputChannels = append(inputChannels, c.ensureChannel(name))
	}

	outputChannel := c.ensureChannel(outputName)

	task := func(ctx context.Context) error {
		return handler(ctx, inputChannels, outputChannel)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputChannel := c.ensureChannel(inputName)

	outputChannels := make([]chan string, 0, len(outputNames))
	for _, name := range outputNames {
		outputChannels = append(outputChannels, c.ensureChannel(name))
	}

	task := func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannels)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, task)
	c.mu.Unlock()
}

func (c *Conveyer) Send(channelName string, data string) error {
	channelRef, ok := c.getChannel(channelName)
	if !ok {
		return ErrChanNotFound
	}

	channelRef <- data

	return nil
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	channelRef, ok := c.getChannel(channelName)
	if !ok {
		return "", ErrChanNotFound
	}

	value, open := <-channelRef
	if !open {
		return "", nil
	}

	return value, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup

	var (
		firstErr error
		once     sync.Once
	)

	for _, handlerFunc := range c.handlers {
		waitGroup.Add(1)

		go func(h func(ctx context.Context) error) {
			defer waitGroup.Done()

			if err := h(ctx); err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(handlerFunc)
	}

	waitGroup.Wait()

	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mu.Unlock()

	return firstErr
}
