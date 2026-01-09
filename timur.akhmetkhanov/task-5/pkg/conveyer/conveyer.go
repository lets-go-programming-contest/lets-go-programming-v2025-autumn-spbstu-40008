package conveyer

import (
	"context"
	"errors"
	"sync"
)

type ConveyerInterface interface {
	RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string)
	RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string)
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
		channels:   make(map[string]chan string),
		bufferSize: size,
		workers:    make([]Worker, 0),
	}
}

func (c *Conveyer) getChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string) {
	inCh := c.getChan(input)
	outCh := c.getChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *Conveyer) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	var inChs []chan string
	for _, name := range inputs {
		inChs = append(inChs, c.getChan(name))
	}
	outCh := c.getChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	})
}

func (c *Conveyer) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	inCh := c.getChan(input)
	var outChs []chan string
	for _, name := range outputs {
		outChs = append(outChs, c.getChan(name))
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.workers))

	for _, w := range c.workers {
		wg.Add(1)
		worker := w
		go func() {
			defer wg.Done()
			if err := worker(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	select {
	case <-ctx.Done():
	case err := <-errCh:
		return err
	}

	c.mu.Lock()
	for _, ch := range c.channels {
		select {
		case <-ch:
		default:
			close(ch)
		}
	}
	c.mu.Unlock()

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return errors.New("chan not found")
	}

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	val, open := <-ch
	if !open {
		return "undefined", nil
	}
	return val, nil
}
