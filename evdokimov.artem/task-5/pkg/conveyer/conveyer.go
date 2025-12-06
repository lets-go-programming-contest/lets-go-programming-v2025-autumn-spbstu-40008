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

type Conveyer struct {
	mutex          sync.RWMutex
	channelsByName map[string]chan string
	handlerList    []func(context.Context) error
	size           int
}

func New(size int) *Conveyer {
	return &Conveyer{
		mutex:          sync.RWMutex{},
		channelsByName: make(map[string]chan string),
		handlerList:    make([]func(context.Context) error, 0),
		size:           size,
	}
}

func (c *Conveyer) ensureChannel(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	channel, exists := c.channelsByName[name]
	if exists {
		return channel
	}

	channel = make(chan string, c.size)
	c.channelsByName[name] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inputChannel := c.ensureChannel(input)
	outputChannel := c.ensureChannel(output)

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputNames []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputNames))

	for i, name := range inputNames {
		inputChannels[i] = c.ensureChannel(name)
	}

	outputChannel := c.ensureChannel(output)

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handler(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	input string,
	outputNames []string,
) {
	inputChannel := c.ensureChannel(input)
	outputChannels := make([]chan string, len(outputNames))

	for i, name := range outputNames {
		outputChannels[i] = c.ensureChannel(name)
	}

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Send(name string, value string) error {
	c.mutex.RLock()
	channel, ok := c.channelsByName[name]
	c.mutex.RUnlock()

	if !ok {
		return ErrChannelNotFound
	}

	channel <- value

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mutex.RLock()
	channel, ok := c.channelsByName[name]
	c.mutex.RUnlock()

	if !ok {
		return "", ErrChannelNotFound
	}

	val, open := <-channel
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
	defer c.mutex.Unlock()

	for _, channel := range c.channelsByName {
		close(channel)
	}

	if err != nil {
		return fmt.Errorf("pipeline failed: %w", err)
	}

	return nil
}
