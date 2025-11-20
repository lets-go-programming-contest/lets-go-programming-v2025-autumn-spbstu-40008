package conveyer

import (
	"context"
	"errors"
	"sync"
)

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
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inChan := c.getOrCreateChan(input)
	outChan := c.getOrCreateChan(output)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inChan, outChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, input []chan string, outputs chan string) error,
	input []string,
	outputs string,
) {
	var inChans []chan string
	for _, name := range input {
		inChans = append(inChans, c.getOrCreateChan(name))
	}

	outChan := c.getOrCreateChan(outputs)
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inChans, outChan)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inChan := c.getOrCreateChan(input)
	var outChans []chan string
	for _, name := range outputs {
		outChans = append(outChans, c.getOrCreateChan(name))
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inChan, outChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var errOnce sync.Once
	var firstErr error

	for _, task := range c.tasks {
		wg.Add(1)

		currentTask := task
		go func() {
			defer wg.Done()
			if err := currentTask(ctx); err != nil {
				errOnce.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}()
	}

	wg.Wait()

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
		return errors.New("chan not found")
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	ch, ok := c.getChan(name)
	if !ok {
		return "", errors.New("chan not found")
	}
	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return data, nil
}
