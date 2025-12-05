package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("chan not found")
)

type conveyer struct {
	size      int
	channels  map[string]chan string
	mu        sync.RWMutex
	handlerMu sync.Mutex

	decorators   []decorator
	multiplexers []multiplexer
	separators   []separator

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type decorator struct {
	fn     func(ctx context.Context, input, output chan string) error
	input  string
	output string
}

type multiplexer struct {
	fn     func(ctx context.Context, inputs []chan string, output chan string) error
	inputs []string
	output string
}

type separator struct {
	fn      func(ctx context.Context, input chan string, outputs []chan string) error
	input   string
	outputs []string
}

func New(size int) *conveyer {
	return &conveyer{
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *conveyer) createChannel(name string) chan string {
	c.mu.RLock()
	if ch, ok := c.channels[name]; ok {
		c.mu.RUnlock()
		return ch
	}
	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.channels[name]
	return ch, ok
}

func (c *conveyer) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error, input, output string) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()
	c.decorators = append(c.decorators, decorator{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()
	c.multiplexers = append(c.multiplexers, multiplexer{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()
	c.separators = append(c.separators, separator{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *conveyer) Send(input string, data string) error {
	ch, ok := c.getChannel(input)
	if !ok {
		return ErrNotFound
	}

	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case ch <- data:
		return nil
	}
}

func (c *conveyer) Recv(output string) (string, error) {
	ch, ok := c.getChannel(output)
	if !ok {
		return "", ErrNotFound
	}

	select {
	case <-c.ctx.Done():
		return "", c.ctx.Err()
	case val, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return val, nil
	}
}

func (c *conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()
	errCh := make(chan error, 1)
	c.startAllHandlers(errCh)
	return c.waitComplete(errCh)
}

func (c *conveyer) startAllHandlers(errCh chan error) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()

	for _, d := range c.decorators {
		c.startHandler(func() error {
			in := c.createChannel(d.input)
			out := c.createChannel(d.output)
			return d.fn(c.ctx, in, out)
		}, errCh)
	}

	for _, m := range c.multiplexers {
		c.startHandler(func() error {
			inputs := make([]chan string, len(m.inputs))
			for i, name := range m.inputs {
				inputs[i] = c.createChannel(name)
			}
			out := c.createChannel(m.output)
			return m.fn(c.ctx, inputs, out)
		}, errCh)
	}

	for _, s := range c.separators {
		c.startHandler(func() error {
			in := c.createChannel(s.input)
			outputs := make([]chan string, len(s.outputs))
			for i, name := range s.outputs {
				outputs[i] = c.createChannel(name)
			}
			return s.fn(c.ctx, in, outputs)
		}, errCh)
	}
}

func (c *conveyer) startHandler(handler func() error, errCh chan error) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		if err := handler(); err != nil {
			select {
			case errCh <- err:
				c.cancel()
			default:
			}
		}
	}()
}

func (c *conveyer) waitComplete(errCh chan error) error {
	go func() {
		c.wg.Wait()
		close(errCh)
	}()

	select {
	case err, ok := <-errCh:
		if ok && err != nil {
			return err
		}
		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
}

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input, output chan string) error, input, output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}
