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
	in := c.ensureChannel(inputName)
	out := c.ensureChannel(outputName)

	task := func(ctx context.Context) error {
		return handler(ctx, in, out)
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
	inChans := make([]chan string, 0, len(inputNames))
	for _, name := range inputNames {
		inChans = append(inChans, c.ensureChannel(name))
	}

	out := c.ensureChannel(outputName)

	task := func(ctx context.Context) error {
		return handler(ctx, inChans, out)
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
	in := c.ensureChannel(inputName)

	outChans := make([]chan string, 0, len(outputNames))
	for _, name := range outputNames {
		outChans = append(outChans, c.ensureChannel(name))
	}

	task := func(ctx context.Context) error {
		return handler(ctx, in, outChans)
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

	val, open := <-channelRef
	if !open {
		return "undefined", nil
	}

	return val, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup
	var firstErr error
	var once sync.Once

	for _, h := range c.handlers {
		waitGroup.Add(1)

		handlerFunc := h

		go func() {
			defer waitGroup.Done()

			if err := handlerFunc(ctx); err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}()
	}

	waitGroup.Wait()

	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mu.Unlock()

	return firstErr
}
