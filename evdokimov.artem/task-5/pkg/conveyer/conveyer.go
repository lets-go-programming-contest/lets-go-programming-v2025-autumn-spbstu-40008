package conveyer

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"sync"
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

	if existing, ok := c.channelsByName[name]; ok {
		return existing
	}

	ch := make(chan string, c.size)
	c.channelsByName[name] = ch

	return ch
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	inputName, outputName string,
) {
	inputCh := c.ensureChannel(inputName)
	outputCh := c.ensureChannel(outputName)

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFunc(ctx, inputCh, outputCh)

	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	inputChs := make([]chan string, len(inputNames))

	for i, name := range inputNames {
		inputChs[i] = c.ensureChannel(name)
	}

	outputCh := c.ensureChannel(outputName)

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChs, outputCh)
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputCh := c.ensureChannel(inputName)
	outputChs := make([]chan string, len(outputNames))

	for i, name := range outputNames {
		outputChs[i] = c.ensureChannel(name)
	}

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFunc(ctx, inputCh, outputChs)
	})
}

func (c *Conveyer) Send(pipe string, value string) error {
	c.mutex.RLock()
	ch, ok := c.channelsByName[pipe]
	c.mutex.RUnlock()

	if !ok {
		return ErrChannelNotFound
	}

	ch <- value

	return nil
}

func (c *Conveyer) Recv(pipe string) (string, error) {
	c.mutex.RLock()
	ch, ok := c.channelsByName[pipe]
	c.mutex.RUnlock()

	if !ok {
		return "", ErrChannelNotFound
	}

	value, open := <-ch
	if !open {
		return UndefinedValue, nil
	}

	return value, nil
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
