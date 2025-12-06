package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mu        sync.RWMutex
	pipelines map[string]chan string
	workers   []func(ctx context.Context) error
	bufSize   int
}

func New(bufferSize int) *Conveyer {
	return &Conveyer{
		pipelines: make(map[string]chan string),
		workers:   make([]func(ctx context.Context) error, 0),
		bufSize:   bufferSize,
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.pipelines[name]; ok {
		return ch
	}

	newCh := make(chan string, c.bufSize)
	c.pipelines[name] = newCh
	return newCh
}

func (c *Conveyer) RegisterDecorator(
	processor func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	input := c.getOrCreateChannel(inputName)
	output := c.getOrCreateChannel(outputName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, input, output)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	processor func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	inputs := make([]chan string, len(inputNames))
	for i, name := range inputNames {
		inputs[i] = c.getOrCreateChannel(name)
	}
	output := c.getOrCreateChannel(outputName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inputs, output)
	})
}

func (c *Conveyer) RegisterSeparator(
	processor func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	input := c.getOrCreateChannel(inputName)
	outputs := make([]chan string, len(outputNames))
	for i, name := range outputNames {
		outputs[i] = c.getOrCreateChannel(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, input, outputs)
	})
}

func (c *Conveyer) Send(pipeName string, data string) error {
	c.mu.RLock()
	ch, exists := c.pipelines[pipeName]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("%w: %s", ErrChanNotFound, pipeName)
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(pipeName string) (string, error) {
	c.mu.RLock()
	ch, exists := c.pipelines[pipeName]
	c.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("%w: %s", ErrChanNotFound, pipeName)
	}

	val, ok := <-ch
	if !ok {
		return "", nil
	}

	return val, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, worker := range c.workers {
		w := worker
		group.Go(func() error {
			return w(ctx)
		})
	}

	return group.Wait()
}