package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mutex      sync.RWMutex
	pipeMap    map[string]chan string
	processors []func(context.Context) error
	bufferSize int
}

func New(size int) *Conveyer {
	return &Conveyer{
		mutex:      sync.RWMutex{},
		pipeMap:    make(map[string]chan string),
		processors: make([]func(context.Context) error, 0),
		bufferSize: size,
	}
}

func (c *Conveyer) getOrInitChannel(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if channel, exists := c.pipeMap[name]; exists {
		return channel
	}

	newChannel := make(chan string, c.bufferSize)
	c.pipeMap[name] = newChannel

	return newChannel
}

func (c *Conveyer) RegisterDecorator(
	action func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	inputChannel := c.getOrInitChannel(inputName)
	outputChannel := c.getOrInitChannel(outputName)

	wrappedTask := func(ctx context.Context) error {
		return action(ctx, inputChannel, outputChannel)
	}

	c.mutex.Lock()
	c.processors = append(c.processors, wrappedTask)
	c.mutex.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	action func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	inputChannels := make([]chan string, 0, len(inputNames))
	for _, name := range inputNames {
		inputChannels = append(inputChannels, c.getOrInitChannel(name))
	}

	outputChannel := c.getOrInitChannel(outputName)

	wrappedTask := func(ctx context.Context) error {
		return action(ctx, inputChannels, outputChannel)
	}

	c.mutex.Lock()
	c.processors = append(c.processors, wrappedTask)
	c.mutex.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	action func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputChannel := c.getOrInitChannel(inputName)

	outputChannels := make([]chan string, 0, len(outputNames))
	for _, name := range outputNames {
		outputChannels = append(outputChannels, c.getOrInitChannel(name))
	}

	wrappedTask := func(ctx context.Context) error {
		return action(ctx, inputChannel, outputChannels)
	}

	c.mutex.Lock()
	c.processors = append(c.processors, wrappedTask)
	c.mutex.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup

	errChan := make(chan error, 1)

	for _, proc := range c.processors {
		waitGroup.Add(1)

		currentProc := proc

		go func(p func(context.Context) error) {
			defer waitGroup.Done()

			if err := p(ctx); err != nil {
				select {
				case errChan <- err:
					cancel()
				default:
				}
			}
		}(currentProc)
	}

	waitGroup.Wait()

	c.mutex.Lock()
	for _, channel := range c.pipeMap {
		close(channel)
	}
	c.mutex.Unlock()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (c *Conveyer) Send(name string, data string) error {
	c.mutex.RLock()
	channel, ok := c.pipeMap[name]
	c.mutex.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mutex.RLock()
	channel, ok := c.pipeMap[name]
	c.mutex.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, isOpen := <-channel
	if !isOpen {
		return "undefined", nil
	}

	return val, nil
}
