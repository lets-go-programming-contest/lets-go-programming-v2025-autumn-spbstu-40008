package conveyer

import (
	"context"
	"errors"
	"sync"
)

type Conveyer struct {
	mu       sync.Mutex
	channels map[string]chan string
	size     int
	tasks    []func(context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
		tasks:    make([]func(context.Context) error, 0),
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

func (c *Conveyer) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string) error {
	inCh := c.getOrCreateChan(input)
	outCh := c.getOrCreateChan(output)
	
	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}
	c.tasks = append(c.tasks, task)
	return nil
}

func (c *Conveyer) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) error {
	var inChs []chan string
	for _, name := range inputs {
		inChs = append(inChs, c.getOrCreateChan(name))
	}
	outCh := c.getOrCreateChan(output)

	task := func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	}
	c.tasks = append(c.tasks, task)
	return nil
}

func (c *Conveyer) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) error {
	inCh := c.getOrCreateChan(input)
	var outChs []chan string
	for _, name := range outputs {
		outChs = append(outChs, c.getOrCreateChan(name))
	}

	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	}
	c.tasks = append(c.tasks, task)
	return nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer c.closeAll()

	errChan := make(chan error, len(c.tasks))
	var wg sync.WaitGroup

	for _, task := range c.tasks {
		wg.Add(1)
		go func(t func(context.Context) error) {
			defer wg.Done()
			err := t(ctx)
			if err != nil {
				errChan <- err
			}
		}(task)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err, ok := <-errChan:
		if !ok {
			return nil
		}
		cancel()
		return err
	}
}

func (c *Conveyer) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, ch := range c.channels {
		select {
		case <-ch:
		default:
		}
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.Lock()
	ch, ok := c.channels[input]
	c.mu.Unlock()
	if !ok {
		return errors.New("chan not found")
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.Lock()
	ch, ok := c.channels[output]
	c.mu.Unlock()
	if !ok {
		return "", errors.New("chan not found")
	}
	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
}