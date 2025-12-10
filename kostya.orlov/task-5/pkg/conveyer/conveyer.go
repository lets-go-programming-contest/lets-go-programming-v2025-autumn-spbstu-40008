package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound                       = errors.New("chan not found")
	ErrChannelExists                      = errors.New("channel already exists")
	ErrUnsupportedDecoratorFunctionType   = errors.New("unsupported decorator function type")
	ErrUnsupportedMultiplexerFunctionType = errors.New("unsupported multiplexer function type")
	ErrUnsupportedSeparatorFunctionType   = errors.New("unsupported separator function type")
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
		tasks:    []func(ctx context.Context) error{},
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

	if chanItem, ok := c.channels[name]; ok {
		return chanItem
	}

	chanItem := make(chan string, c.size)
	c.channels[name] = chanItem

	return chanItem
}

func (c *Conveyer) getChan(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	chanItem, ok := c.channels[name]

	return chanItem, ok
}

func (c *Conveyer) RegisterDecorator(action interface{}, input, output string) {
	inChan := c.getOrCreateChan(input)
	outChan := c.getOrCreateChan(output)

	var task func(context.Context) error

	switch funcType := action.(type) {
	case func(context.Context, chan string, chan string) error:
		task = func(ctx context.Context) error {
			return funcType(ctx, inChan, outChan)
		}
	case func(context.Context, chan string, []chan string) error:
		task = func(ctx context.Context) error {
			return funcType(ctx, inChan, []chan string{outChan})
		}
	default:
		task = func(context.Context) error {
			return ErrUnsupportedDecoratorFunctionType
		}
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(action interface{}, input []string, outputs string) {
	inChans := make([]chan string, 0, len(input))
	for _, name := range input {
		inChans = append(inChans, c.getOrCreateChan(name))
	}

	outChan := c.getOrCreateChan(outputs)

	var task func(context.Context) error

	switch funcType := action.(type) {
	case func(context.Context, []chan string, chan string) error:
		task = func(ctx context.Context) error {
			return funcType(ctx, inChans, outChan)
		}
	case func(context.Context, []chan string, []chan string) error:
		task = func(ctx context.Context) error {
			return funcType(ctx, inChans, []chan string{outChan})
		}
	default:
		task = func(context.Context) error {
			return ErrUnsupportedMultiplexerFunctionType
		}
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(action interface{}, input string, outputs []string) {
	inChan := c.getOrCreateChan(input)

	outChans := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outChans = append(outChans, c.getOrCreateChan(name))
	}

	var task func(context.Context) error

	switch funcType := action.(type) {
	case func(context.Context, chan string, []chan string) error:
		task = func(ctx context.Context) error {
			return funcType(ctx, inChan, outChans)
		}
	default:
		task = func(context.Context) error {
			return ErrUnsupportedSeparatorFunctionType
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

	var once sync.Once

	var firstErr error

	for i := range c.tasks {
		waitGroup.Add(1)

		currentTask := c.tasks[i]

		go func() {
			defer waitGroup.Done()

			if err := currentTask(ctx); err != nil {
				once.Do(func() {
					firstErr = err

					cancel()
				})
			}
		}()
	}

	waitGroup.Wait()

	c.mu.Lock()
	for _, chanItem := range c.channels {
		func(channel chan string) {
			defer func() {
				_ = recover()
			}()
			close(channel)
		}(chanItem)
	}
	c.mu.Unlock()

	return firstErr
}

func (c *Conveyer) Send(name, data string) error {
	chanItem, ok := c.getChan(name)
	if !ok {
		return ErrChanNotFound
	}

	chanItem <- data

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	chanItem, exists := c.getChan(name)
	if !exists {
		return "", ErrChanNotFound
	}

	data, isOpen := <-chanItem
	if !isOpen {
		return "undefined", nil
	}

	return data, nil
}
