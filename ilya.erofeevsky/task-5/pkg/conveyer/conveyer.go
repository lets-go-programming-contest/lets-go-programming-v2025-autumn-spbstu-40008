package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("chan not found")

const undefinedValue = "undefined"

type Conveyer struct {
	size           int
	channelsByName map[string]chan string
	mutex          sync.RWMutex
	handlerList    []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:           size,
		channelsByName: make(map[string]chan string),
		handlerList:    []func(context.Context) error{},
		mutex:          sync.RWMutex{},
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

func (c *Conveyer) getChannel(name string) chan string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.channelsByName[name]
}

func (c *Conveyer) RegisterDecorator(
	handlerFn func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inputCh := c.ensureChannel(input)
	outputCh := c.ensureChannel(output)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFn(ctx, inputCh, outputCh)
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputCh := c.ensureChannel(input)
	outputsChs := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputsChs[i] = c.ensureChannel(name)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFn(ctx, inputCh, outputsChs)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputsChs := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputsChs[i] = c.ensureChannel(name)
	}

	outputCh := c.ensureChannel(output)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.handlerList = append(c.handlerList, func(ctx context.Context) error {
		return handlerFn(ctx, inputsChs, outputCh)
	})
}

func (c *Conveyer) Send(input, data string) error {
	inputCh := c.getChannel(input)
	if inputCh == nil {
		return ErrChannelNotFound
	}
	inputCh <- data
	
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	outputCh := c.getChannel(output)
	if outputCh == nil {
		return "", ErrChannelNotFound
	}

	val, ok := <-outputCh
	if !ok {
		return undefinedValue, nil
	}

	return val, nil
}

func (c *Conveyer) Run(executionContext context.Context) error {
	defer func() {
		c.mutex.RLock()
		defer c.mutex.RUnlock()
		
		for _, channel := range c.channelsByName { 
			select {
			case <-executionContext.Done():
			default:
				close(channel)
			}
		}
	}()

	errorGroup, operationContext := errgroup.WithContext(executionContext)

	c.mutex.RLock()

	for _, handler := range c.handlerList {
		
		errorGroup.Go(func() error {
			return handler(operationContext)
		})
	}
	c.mutex.RUnlock()

	if runError := errorGroup.Wait(); runError != nil {
		return fmt.Errorf("run pipeline: %w", runError)
	}
	
	return nil
}
