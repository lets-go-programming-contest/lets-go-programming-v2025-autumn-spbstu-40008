package conveyer

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type handlerType int

const (
	decoratorType handlerType = iota
	multiplexerType
	separatorType
)

type handlerFunc struct {
	fn      interface{}
	inputs  []string
	outputs []string
	htype   handlerType
}

type Conveyer struct {
	size     int
	channels map[string]chan string
	mu       sync.RWMutex
	handlers []handlerFunc
	cancel   context.CancelFunc
	closed   atomic.Bool
	wg       sync.WaitGroup
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input, output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChannel(input)
	c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, handlerFunc{
		fn:      fn,
		inputs:  []string{input},
		outputs: []string{output},
		htype:   decoratorType,
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, in := range inputs {
		c.getOrCreateChannel(in)
	}
	c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, handlerFunc{
		fn:      fn,
		inputs:  inputs,
		outputs: []string{output},
		htype:   multiplexerType,
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChannel(input)
	for _, out := range outputs {
		c.getOrCreateChannel(out)
	}

	c.handlers = append(c.handlers, handlerFunc{
		fn:      fn,
		inputs:  []string{input},
		outputs: outputs,
		htype:   separatorType,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	c.wg.Add(len(c.handlers))
	errCh := make(chan error, len(c.handlers))

	for _, h := range c.handlers {
		go func(h handlerFunc) {
			defer c.wg.Done()

			c.mu.RLock()
			inputs := make([]chan string, len(h.inputs))
			for i, name := range h.inputs {
				inputs[i] = c.channels[name]
			}
			outputs := make([]chan string, len(h.outputs))
			for i, name := range h.outputs {
				outputs[i] = c.channels[name]
			}
			c.mu.RUnlock()

			var err error
			switch h.htype {
			case decoratorType:
				err = h.fn.(func(context.Context, chan string, chan string) error)(ctx, inputs[0], outputs[0])
			case multiplexerType:
				err = h.fn.(func(context.Context, []chan string, chan string) error)(ctx, inputs, outputs[0])
			case separatorType:
				err = h.fn.(func(context.Context, chan string, []chan string) error)(ctx, inputs[0], outputs)
			}

			if err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(h)
	}

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		cancel()
		c.wg.Wait()
		c.closeAllChannels()
		return ctx.Err()

	case err := <-errCh:
		cancel()
		c.wg.Wait()
		c.closeAllChannels()
		return err

	case <-done:
		cancel()
		c.closeAllChannels()
		return nil
	}
}

func (c *Conveyer) closeAllChannels() {
	if c.closed.Load() {
		return
	}
	c.closed.Store(true)

	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
}

func (c *Conveyer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	c.closeAllChannels()
}

func (c *Conveyer) Send(input string, data string) error {
	if c.closed.Load() {
		return errors.New("conveyer closed")
	}

	c.mu.RLock()
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return errors.New("chan not found")
	}

	select {
	case ch <- data:
		return nil
	default:
		return errors.New("send timeout")
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	if c.closed.Load() {
		return "", errors.New("conveyer closed")
	}

	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
