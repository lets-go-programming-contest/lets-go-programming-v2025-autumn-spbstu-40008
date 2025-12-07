package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
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
	channels   map[string]chan string
	mu         sync.Mutex
	decorators []struct {
		fn     ModifierFunc
		input  string
		output string
	}
	multiplexers []struct {
		fn     MultiplexerFunc
		inputs []string
		output string
	}
	separators []struct {
		fn      SeparatorFunc
		input   string
		outputs []string
	}
}

var ErrChanNotFound = errors.New("chan not found")

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		channels: make(map[string]chan string),
		mu:       sync.Mutex{},
		decorators: []struct {
			fn            ModifierFunc
			input, output string
		}{},
		multiplexers: []struct {
			fn     MultiplexerFunc
			inputs []string
			output string
		}{},
		separators: []struct {
			fn      SeparatorFunc
			input   string
			outputs []string
		}{},
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

func (c *conveyerImpl) RegisterDecorator(fn ModifierFunc, input string, output string) {
	c.decorators = append(c.decorators, struct {
		fn     ModifierFunc
		input  string
		output string
	}{fn: fn, input: input, output: output})
}

func (c *conveyerImpl) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
	c.multiplexers = append(c.multiplexers, struct {
		fn     MultiplexerFunc
		inputs []string
		output string
	}{fn: fn, inputs: inputs, output: output})
}

func (c *conveyerImpl) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
	c.separators = append(c.separators, struct {
		fn      SeparatorFunc
		input   string
		outputs []string
	}{fn: fn, input: input, outputs: outputs})
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

	g, gCtx := errgroup.WithContext(ctx)

	for _, dec := range c.decorators {
		dec := dec
		g.Go(func() error {
			inCh, _ := c.getChannel(dec.input)
			outCh, _ := c.getChannel(dec.output)
			return dec.fn(gCtx, inCh, outCh)
		})
	}

	for _, mux := range c.multiplexers {
		mux := mux
		g.Go(func() error {
			var inChs []chan string
			for _, inName := range mux.inputs {
				ch, _ := c.getChannel(inName)
				inChs = append(inChs, ch)
			}
			outCh, _ := c.getChannel(mux.output)
			return mux.fn(gCtx, inChs, outCh)
		})
	}

	for _, sep := range c.separators {
		sep := sep
		g.Go(func() error {
			inCh, _ := c.getChannel(sep.input)
			var outChs []chan string
			for _, outName := range sep.outputs {
				ch, _ := c.getChannel(outName)
				outChs = append(outChs, ch)
			}
			return sep.fn(gCtx, inCh, outChs)
		})
	}

	if err := g.Wait(); err != nil {
		c.mu.Lock()
		for _, ch := range c.channels {
			close(ch)
		}
		c.mu.Unlock()
		return err
	}

	c.mu.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mu.Unlock()
	return nil
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, exists := c.getChannel(input)
	if !exists {
		return ErrChanNotFound
	}

	select {
	case ch <- data:
		return nil
	case <-context.Background().Done():
		return context.Canceled
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)
	if !exists {
		return "", ErrChanNotFound
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
