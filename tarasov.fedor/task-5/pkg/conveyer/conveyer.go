package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channels map[string]chan string
	tasks    []func(ctx context.Context) error
	size     int
	mu       sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		tasks:    make([]func(ctx context.Context) error, 0),
		size:     size,
		mu:       sync.RWMutex{},
	}
}

func (c *Conveyer) getOrCreateChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) getChan(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[name]

	return ch, ok
}

func (c *Conveyer) RegisterDecorator(
	action func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inChan := c.getOrCreateChan(input)
	outChan := c.getOrCreateChan(output)

	task := func(ctx context.Context) error {
		return action(ctx, inChan, outChan)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	action func(ctx context.Context, input []chan string, outputs chan string) error,
	input []string,
	outputs string,
) {
	inChans := make([]chan string, 0, len(input))
	for _, name := range input {
		inChans = append(inChans, c.getOrCreateChan(name))
	}

	outChan := c.getOrCreateChan(outputs)

	task := func(ctx context.Context) error {
		return action(ctx, inChans, outChan)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	action func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inChan := c.getOrCreateChan(input)

	outChans := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outChans = append(outChans, c.getOrCreateChan(name))
	}

	task := func(ctx context.Context) error {
		return action(ctx, inChan, outChans)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup

	var errOnce sync.Once

	var firstErr error

	for _, task := range c.tasks {
		waitGroup.Add(1)

		currentTask := task

		go func() {
			defer waitGroup.Done()

			if err := currentTask(ctx); err != nil {
				errOnce.Do(func() {
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

func (c *Conveyer) Send(name string, data string) error {
	ch, ok := c.getChan(name)
	if !ok {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	ch, exists := c.getChan(name)
	if !exists {
		return "", ErrChanNotFound
	}

	data, isOpen := <-ch

	if !isOpen {
		return "undefined", nil
	}

	return data, nil
}
