package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("channel not found")

const UndefinedValue = "undefined"

type Conveyer struct {
	mutex          sync.RWMutex
	channelsByName map[string]chan string
	handlerList    []func(context.Context) error
	size           int
}

func New(size int) *Conveyer {
	return &Conveyer{
		channelsByName: make(map[string]chan string),
		handlerList:    make([]func(context.Context) error, 0),
		size:           size,
	}
}

func (c *Conveyer) ensureChannel(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ch, ok := c.channelsByName[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channelsByName[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	inputChannelName string,
	outputChannelName string,
) {
	inputChannel := c.ensureChannel(inputChannelName)
	outputChannel := c.ensureChannel(outputChannelName)

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return fn(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputChannelNames []string,
	outputChannelName string,
) {
	inputChannels := make([]chan string, len(inputChannelNames))
	for i, name := range inputChannelNames {
		inputChannels[i] = c.ensureChannel(name)
	}

	outputChannel := c.ensureChannel(outputChannelName)

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return fn(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	inputChannelName string,
	outputChannelNames []string,
) {
	inputChannel := c.ensureChannel(inputChannelName)

	outputChannels := make([]chan string, len(outputChannelNames))
	for i, name := range outputChannelNames {
		outputChannels[i] = c.ensureChannel(name)
	}

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return fn(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Send(channelName string, value string) error {
	c.mutex.RLock()
	ch, ok := c.channelsByName[channelName]
	c.mutex.RUnlock()

	if !ok {
		return ErrChannelNotFound
	}

	ch <- value
	return nil
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	c.mutex.RLock()
	ch, ok := c.channelsByName[channelName]
	c.mutex.RUnlock()

	if !ok {
		return "", ErrChannelNotFound
	}

	val, open := <-ch
	if !open {
		return UndefinedValue, nil
	}

	return val, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlerList {
		h := handler
		group.Go(func() error {
			return h(ctx)
		})
	}

	err := group.Wait()

	c.mutex.Lock()
	for _, ch := range c.channelsByName {
		close(ch)
	}
	c.mutex.Unlock()

	return err
}
