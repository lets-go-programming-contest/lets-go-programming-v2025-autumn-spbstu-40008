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

	if channel, ok := c.channelsByName[name]; ok {
		return channel
	}

	channel := make(chan string, c.size)
	c.channelsByName[name] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	inputChannel := c.ensureChannel(inputName)
	outputChannel := c.ensureChannel(outputName)

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	inputChannels := make([]chan string, len(inputNames))

	for i, name := range inputNames {
		inputChannels[i] = c.ensureChannel(name)
	}

	outputChannel := c.ensureChannel(outputName)

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputChannel := c.ensureChannel(inputName)
	outputChannels := make([]chan string, len(outputNames))

	for i, name := range outputNames {
		outputChannels[i] = c.ensureChannel(name)
	}

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Send(pipe string, value string) error {
	c.mutex.RLock()
	channel, ok := c.channelsByName[pipe]
	c.mutex.RUnlock()

	if !ok {
		return ErrChannelNotFound
	}

	channel <- value

	return nil
}

func (c *Conveyer) Recv(pipe string) (string, error) {
	c.mutex.RLock()
	channel, ok := c.channelsByName[pipe]
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

func (c *Conveyer) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlerList {
		handler := h

		grp.Go(func() error {
			return handler(ctx)
		})
	}

	err := grp.Wait()

	c.mutex.Lock()
	for _, channel := range c.channelsByName {
		close(channel)
	}
	c.mutex.Unlock()

	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}
