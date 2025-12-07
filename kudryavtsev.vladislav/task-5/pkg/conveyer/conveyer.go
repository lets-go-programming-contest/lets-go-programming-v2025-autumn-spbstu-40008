package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mutex sync.Mutex

	channels map[string]chan string

	tasks []func(context.Context) error

	size int
}

func New(size int) *Conveyer {

	return &Conveyer{

		mutex: sync.Mutex{},

		channels: make(map[string]chan string),

		tasks: make([]func(context.Context) error, 0),

		size: size,
	}

}

func (c *Conveyer) RegisterDecorator(

	decorator func(context.Context, chan string, chan string) error,

	input string,

	output string,

) error {

	c.mutex.Lock()

	if _, exists := c.channels[input]; !exists {

		c.channels[input] = make(chan string, c.size)

	}

	inCh := c.channels[input]

	if _, exists := c.channels[output]; !exists {

		c.channels[output] = make(chan string, c.size)

	}

	outCh := c.channels[output]

	c.mutex.Unlock()

	task := func(ctx context.Context) error {

		return decorator(ctx, inCh, outCh)

	}

	c.tasks = append(c.tasks, task)

	return nil

}

func (c *Conveyer) RegisterMultiplexer(

	multiplexer func(context.Context, []chan string, chan string) error,

	inputs []string,

	output string,

) error {

	inChs := make([]chan string, 0, len(inputs))

	c.mutex.Lock()

	for _, name := range inputs {

		if _, exists := c.channels[name]; !exists {

			c.channels[name] = make(chan string, c.size)

		}

		inChs = append(inChs, c.channels[name])

	}

	if _, exists := c.channels[output]; !exists {

		c.channels[output] = make(chan string, c.size)

	}

	outCh := c.channels[output]

	c.mutex.Unlock()

	task := func(ctx context.Context) error {

		return multiplexer(ctx, inChs, outCh)

	}

	c.tasks = append(c.tasks, task)

	return nil

}

func (c *Conveyer) RegisterSeparator(

	separator func(context.Context, chan string, []chan string) error,

	input string,

	outputs []string,

) error {

	c.mutex.Lock()

	if _, exists := c.channels[input]; !exists {

		c.channels[input] = make(chan string, c.size)

	}

	inCh := c.channels[input]

	outChs := make([]chan string, 0, len(outputs))

	for _, name := range outputs {

		if _, exists := c.channels[name]; !exists {

			c.channels[name] = make(chan string, c.size)

		}

		outChs = append(outChs, c.channels[name])

	}

	c.mutex.Unlock()

	task := func(ctx context.Context) error {

		return separator(ctx, inCh, outChs)

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

		go func(work func(context.Context) error) {

			defer waitGroup.Done()

			err := work(ctxCancel)

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

	for _, channel := range c.channels {

		func() {

			defer func() {

				_ = recover()

			}()

			close(channel)

		}()

	}

}

func (c *Conveyer) Send(input string, data string) error {

	c.mutex.Lock()

	targetChan, exists := c.channels[input]

	c.mutex.Unlock()

	if !exists {

		return ErrChanNotFound

	}

	targetChan <- data

	return nil

}

func (c *Conveyer) Recv(output string) (string, error) {

	c.mutex.Lock()

	targetChan, exists := c.channels[output]

	c.mutex.Unlock()

	if !exists {

		return "", ErrChanNotFound

	}

	val, isOpen := <-targetChan

	if !isOpen {

		return "undefined", nil

	}

	return val, nil

}
