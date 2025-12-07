package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound  = errors.New("chan not found")
	ErrChannelExists = errors.New("channel already exists")
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
		mu:       sync.RWMutex{},
	}
}

func (c *Conveyer) AddChannel(channelID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.channels[channelID]; exists {
		return ErrChannelExists
	}

	c.channels[channelID] = make(chan string, c.size)
	return nil
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
	action interface{},
	input string,
	output string,
) {
	inChan := c.getOrCreateChan(input)
	outChan := c.getOrCreateChan(output)

	task := func(ctx context.Context) error {
		switch fn := action.(type) {
		case func(context.Context, chan string, chan string) error:
			return fn(ctx, inChan, outChan)
		case func(context.Context, chan string, []chan string) error:
			return fn(ctx, inChan, []chan string{outChan})
		default:
			return errors.New("unsupported decorator function type")
		}
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	action interface{},
	input []string,
	outputs string,
) {
	inChans := make([]chan string, 0, len(input))
	for _, name := range input {
		inChans = append(inChans, c.getOrCreateChan(name))
	}

	outChan := c.getOrCreateChan(outputs)

	task := func(ctx context.Context) error {
		switch fn := action.(type) {
		case func(context.Context, []chan string, chan string) error:
			return fn(ctx, inChans, outChan)
		case func(context.Context, []chan string, []chan string) error:
			return fn(ctx, inChans, []chan string{outChan})
		default:
			return errors.New("unsupported multiplexer function type")
		}
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	action interface{},
	input string,
	outputs []string,
) {
	inChan := c.getOrCreateChan(input)

	outChans := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outChans = append(outChans, c.getOrCreateChan(name))
	}

	task := func(ctx context.Context) error {
		switch fn := action.(type) {
		case func(context.Context, chan string, []chan string) error:
			return fn(ctx, inChan, outChans)
		default:
			return errors.New("unsupported separator function type")
		}
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
		func(channel chan string) {
			defer func() {
				if r := recover(); r != nil {
				
				}
			}()
			close(channel)
		}(ch)
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