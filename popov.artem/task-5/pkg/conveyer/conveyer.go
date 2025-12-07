package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type conveyerImpl struct {
	size          int
	channels      map[string]chan string
	decorators    []struct{ fn DecoratorFunc; input, output string }
	multiplexers  []struct{ fn MultiplexerFunc; inputs []string; output string }
	separators    []struct{ fn SeparatorFunc; input string; outputs []string }
	mu            sync.Mutex
}

type conveyer interface {
	RegisterDecorator(fn DecoratorFunc, input, output string)
	RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string)
	RegisterSeparator(fn SeparatorFunc, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

func New(size int) conveyer {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *conveyerImpl) getChannel(name string) (chan string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch, exists := c.channels[name]
	if !exists {
		ch = make(chan string, c.size)
		c.channels[name] = ch
	}
	return ch, exists
}

func (c *conveyerImpl) RegisterDecorator(fn DecoratorFunc, input, output string) {
	c.decorators = append(c.decorators, struct {
		fn     DecoratorFunc
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var runWg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, dec := range c.decorators {
		runWg.Add(1)
		go func(d struct {
			fn     DecoratorFunc
			input  string
			output string
		}) {
			defer runWg.Done()
			inCh, _ := c.getChannel(d.input)
			outCh, _ := c.getChannel(d.output)
			if err := d.fn(ctx, inCh, outCh); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(dec)
	}

	for _, mux := range c.multiplexers {
		runWg.Add(1)
		go func(m struct {
			fn     MultiplexerFunc
			inputs []string
			output string
		}) {
			defer runWg.Done()
			var inChs []chan string
			for _, inName := range m.inputs {
				ch, _ := c.getChannel(inName)
				inChs = append(inChs, ch)
			}
			outCh, _ := c.getChannel(m.output)
			if err := m.fn(ctx, inChs, outCh); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(mux)
	}

	for _, sep := range c.separators {
		runWg.Add(1)
		go func(s struct {
			fn      SeparatorFunc
			input   string
			outputs []string
		}) {
			defer runWg.Done()
			inCh, _ := c.getChannel(s.input)
			var outChs []chan string
			for _, outName := range s.outputs {
				ch, _ := c.getChannel(outName)
				outChs = append(outChs, ch)
			}
			if err := s.fn(ctx, inCh, outChs); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(sep)
	}

	done := make(chan struct{})
	go func() {
		runWg.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		c.mu.Lock()
		for _, ch := range c.channels {
			close(ch)
		}
		c.mu.Unlock()
		return err
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
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)
	if !exists {
		return "", ErrChanNotFound
	}
	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return data, nil
}
