// pkg/conveyer/conveyer.go
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

	if pipe, exists := c.pipeMap[name]; exists {
		return pipe
	}

	newPipe := make(chan string, c.bufferSize)
	c.pipeMap[name] = newPipe

	return newPipe
}

func (c *Conveyer) RegisterDecorator(
	action func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	inputPipe := c.getOrInitChannel(inputName)
	outputPipe := c.getOrInitChannel(outputName)

	wrappedTask := func(ctx context.Context) error {
		return action(ctx, inputPipe, outputPipe)
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
	inputPipes := make([]chan string, 0, len(inputNames))
	for _, name := range inputNames {
		inputPipes = append(inputPipes, c.getOrInitChannel(name))
	}

	outputPipe := c.getOrInitChannel(outputName)

	wrappedTask := func(ctx context.Context) error {
		return action(ctx, inputPipes, outputPipe)
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
	inputPipe := c.getOrInitChannel(inputName)

	outputPipes := make([]chan string, 0, len(outputNames))
	for _, name := range outputNames {
		outputPipes = append(outputPipes, c.getOrInitChannel(name))
	}

	wrappedTask := func(ctx context.Context) error {
		return action(ctx, inputPipe, outputPipes)
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

	// runProcessor declared separately to satisfy wsl/go-linter style rules.
	runProcessor := func(proc func(context.Context) error) {
		defer waitGroup.Done()

		if err := proc(ctx); err != nil {
			select {
			case errChan <- err:
				cancel()
			default:
			}
		}
	}

	for _, proc := range c.processors {
		waitGroup.Add(1)
		go runProcessor(proc)
	}

	waitGroup.Wait()

	c.mutex.Lock()
	for _, pipe := range c.pipeMap {
		close(pipe)
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
	pipe, ok := c.pipeMap[name]
	c.mutex.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	pipe <- data

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mutex.RLock()
	pipe, ok := c.pipeMap[name]
	c.mutex.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, isOpen := <-pipe
	if !isOpen {
		return "undefined", nil
	}

	return val, nil
}
