package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type ConveyerInterface interface {
	RegisterDecorator(
		decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type Worker func(ctx context.Context) error

type Conveyer struct {
	mu         sync.RWMutex
	channels   map[string]chan string
	bufferSize int
	workers    []Worker
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:         sync.RWMutex{},
		channels:   make(map[string]chan string),
		bufferSize: size,
		workers:    make([]Worker, 0),
	}
}

func (c *Conveyer) getChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.bufferSize)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChannel := c.getChan(input)
	outputChannel := c.getChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		inputChannels = append(inputChannels, c.getChan(name))
	}

	outputChannel := c.getChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChannel := c.getChan(input)
	outputChannels := make([]chan string, 0, len(outputs))

	for _, name := range outputs {
		outputChannels = append(outputChannels, c.getChan(name))
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup

	errCh := make(chan error, len(c.workers))

	for _, workerFunc := range c.workers {
		waitGroup.Add(1)

		worker := workerFunc

		go func() {
			defer waitGroup.Done()

			if err := worker(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	doneCh := make(chan struct{})
	go func() {
		waitGroup.Wait()
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		return err
	case <-doneCh:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	channel, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, open := <-channel
	if !open {
		return "undefined", nil
	}

	return val, nil
}
