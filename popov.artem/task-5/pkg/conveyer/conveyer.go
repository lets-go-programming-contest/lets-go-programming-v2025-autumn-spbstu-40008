package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	stages   []func(context.Context) error
	bufSize  int
}

func New(size int) *conveyer {
	return &conveyer{
		channels: make(map[string]chan string),
		stages:   make([]func(context.Context) error, 0),
		bufSize:  size,
	}
}

func (c *conveyer) lookupChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.bufSize)
	c.channels[name] = ch
	return ch
}

func (c *conveyer) RegisterDecorator(fn func(context.Context, chan string, chan string) error, in, out string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	input := c.lookupChan(in)
	output := c.lookupChan(out)
	c.stages = append(c.stages, func(ctx context.Context) error {
		return fn(ctx, input, output)
	})
}

func (c *conveyer) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, ins []string, out string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	inputs := make([]chan string, len(ins))
	for i, name := range ins {
		inputs[i] = c.lookupChan(name)
	}
	output := c.lookupChan(out)
	c.stages = append(c.stages, func(ctx context.Context) error {
		return fn(ctx, inputs, output)
	})
}

func (c *conveyer) RegisterSeparator(fn func(context.Context, chan string, []chan string) error, in string, outs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	input := c.lookupChan(in)
	outputs := make([]chan string, len(outs))
	for i, name := range outs {
		outputs[i] = c.lookupChan(name)
	}
	c.stages = append(c.stages, func(ctx context.Context) error {
		return fn(ctx, input, outputs)
	})
}

func (c *conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for _, stage := range c.stages {
		wg.Add(1)
		go func(s func(context.Context) error) {
			defer wg.Done()
			if err := s(ctx); err != nil {
				select {
				case errCh <- err:
					cancel()
				default:
				}
			}
		}(stage)
	}

	wg.Wait()

	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mu.Unlock()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (c *conveyer) Send(chName string, data string) error {
	c.mu.RLock()
	ch, exists := c.channels[chName]
	c.mu.RUnlock()
	if !exists {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *conveyer) Recv(chName string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[chName]
	c.mu.RUnlock()
	if !exists {
		return "", ErrChanNotFound
	}
	msg, open := <-ch
	if !open {
		return "undefined", nil
	}
	return msg, nil
}
