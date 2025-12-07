package conveyer

import (
	"context"
	"errors"
	"sync"
)

type ModifierFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type Conveyer interface {
	RegisterDecorator(fn ModifierFunc, input string, output string)
	RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string)
	RegisterSeparator(fn SeparatorFunc, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyerImpl struct {
	channels     map[string]chan string
	mu           sync.Mutex
	decorators   []struct{ fn ModifierFunc; input, output string }
	multiplexers []struct{ fn MultiplexerFunc; inputs []string; output string }
	separators   []struct{ fn SeparatorFunc; input string; outputs []string }
}

func New(size int) Conveyer {
	return &conveyerImpl{
		channels: make(map[string]chan string),
	}
}

func (c *conveyerImpl) getOrCreateChannel(name string, size int) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, exists := c.channels[name]; exists {
		return ch
	}
	ch := make(chan string, size)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) getChannel(name string) (chan string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch, exists := c.channels[name]
	return ch, exists
}

func (c *conveyerImpl) RegisterDecorator(fn ModifierFunc, input, output string) {
	c.decorators = append(c.decorators, struct{ fn ModifierFunc; input, output string }{fn: fn, input: input, output: output})
}

func (c *conveyerImpl) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	c.multiplexers = append(c.multiplexers, struct{ fn MultiplexerFunc; inputs []string; output string }{fn: fn, inputs: inputs, output: output})
}

func (c *conveyerImpl) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	c.separators = append(c.separators, struct{ fn SeparatorFunc; input string; outputs []string }{fn: fn, input: input, outputs: outputs})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	for _, dec := range c.decorators {
		c.getOrCreateChannel(dec.input, cap(c.getOrCreateChannel(dec.input, 0)))
		c.getOrCreateChannel(dec.output, cap(c.getOrCreateChannel(dec.output, 0)))
	}
	for _, mux := range c.multiplexers {
		for _, in := range mux.inputs {
			c.getOrCreateChannel(in, cap(c.getOrCreateChannel(in, 0)))
		}
		c.getOrCreateChannel(mux.output, cap(c.getOrCreateChannel(mux.output, 0)))
	}
	for _, sep := range c.separators {
		c.getOrCreateChannel(sep.input, cap(c.getOrCreateChannel(sep.input, 0)))
		for _, out := range sep.outputs {
			c.getOrCreateChannel(out, cap(c.getOrCreateChannel(out, 0)))
		}
	}

	var wg sync.WaitGroup

	for _, dec := range c.decorators {
		wg.Add(1)
		go func(d struct{ fn ModifierFunc; input, output string }) {
			defer wg.Done()
			inCh, _ := c.getChannel(d.input)
			outCh, _ := c.getChannel(d.output)
			d.fn(ctx, inCh, outCh)
		}(dec)
	}

	for _, mux := range c.multiplexers {
		wg.Add(1)
		go func(m struct{ fn MultiplexerFunc; inputs []string; output string }) {
			defer wg.Done()
			var inChs []chan string
			for _, inName := range m.inputs {
				ch, _ := c.getChannel(inName)
				inChs = append(inChs, ch)
			}
			outCh, _ := c.getChannel(m.output)
			m.fn(ctx, inChs, outCh)
		}(mux)
	}

	for _, sep := range c.separators {
		wg.Add(1)
		go func(s struct{ fn SeparatorFunc; input string; outputs []string }) {
			defer wg.Done()
			inCh, _ := c.getChannel(s.input)
			var outChs []chan string
			for _, outName := range s.outputs {
				ch, _ := c.getChannel(outName)
				outChs = append(outChs, ch)
			}
			s.fn(ctx, inCh, outChs)
		}(sep)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		c.mu.Lock()
		for _, ch := range c.channels {
			close(ch)
		}
		c.mu.Unlock()
		return nil
	case <-ctx.Done():
		c.mu.Lock()
		for _, ch := range c.channels {
			close(ch)
		}
		c.mu.Unlock()
		return ctx.Err()
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, exists := c.getChannel(input)
	if !exists {
		return errors.New("chan not found")
	}
	select {
	case ch <- 
		return nil
	case <-context.Background().Done():
		return context.Canceled
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)
	if !exists {
		return "", errors.New("chan not found")
	}
	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	case <-context.Background().Done():
		return "", context.Canceled
	}
}