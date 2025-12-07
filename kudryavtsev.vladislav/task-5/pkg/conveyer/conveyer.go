package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
)

type Conveyer struct {
	mutex    sync.Mutex
	channels map[string]chan string
	tasks    []func(context.Context) error
	size     int
}

func New(size int) *Conveyer {
	return &Conveyer{
		mutex:    sync.Mutex{},
		channels: make(map[string]chan string),
		tasks:    make([]func(context.Context) error, 0),
		size:     size,
	}
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) error {
	c.mutex.Lock()

	if _, ok := c.channels[input]; !ok {
		c.channels[input] = make(chan string, c.size)
	}
	inCh := c.channels[input]

	if _, ok := c.channels[output]; !ok {
		c.channels[output] = make(chan string, c.size)
	}
	outCh := c.channels[output]

	c.mutex.Unlock()

	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}

	c.tasks = append(c.tasks, task)

	return nil
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) error {
	inChs := make([]chan string, 0, len(inputs))

	c.mutex.Lock()

	for _, name := range inputs {
		if _, ok := c.channels[name]; !ok {
			c.channels[name] = make(chan string, c.size)
		}

		inChs = append(inChs, c.channels[name])
	}

	if _, ok := c.channels[output]; !ok {
		c.channels[output] = make(chan string, c.size)
	}
	outCh := c.channels[output]

	c.mutex.Unlock()

	task := func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	}

	c.tasks = append(c.tasks, task)

	return nil
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) error {
	c.mutex.Lock()

	if _, ok := c.channels[input]; !ok {
		c.channels[input] = make(chan string, c.size)
	}
	inCh := c.channels[input]

	outChs := make([]chan string, 0, len(outputs))

	for _, name := range outputs {
		if _, ok := c.channels[name]; !ok {
			c.channels[name] = make(chan string, c.size)
		}

		outChs = append(outChs, c.channels[name])
	}

	c.mutex.Unlock()

	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	}

	c.tasks = append(c.tasks, task)

	return nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup

	errChan := make(chan error, len(c.tasks))

	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, task := range c.tasks {
		waitGroup.Add(1)

		go func(tsk func(context.Context) error) {
			defer waitGroup.Done()

			err := tsk(ctxCancel)
			if err != nil {
				errChan <- err
			}
		}(task)
	}

	done := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	var err error

	select {
	case err = <-errChan:
		cancel()
	case <-ctx.Done():
	case <-done:
	}

	c.closeAllChannels()

	return err
}

func (c *Conveyer) closeAllChannels() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, ch := range c.channels {
		func() {
			defer func() {
				_ = recover()
			}()

			close(ch)
		}()
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mutex.Lock()
	ch, ok := c.channels[input]
	c.mutex.Unlock()

	if !ok {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mutex.Lock()
	ch, ok := c.channels[output]
	c.mutex.Unlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return val, nil
}