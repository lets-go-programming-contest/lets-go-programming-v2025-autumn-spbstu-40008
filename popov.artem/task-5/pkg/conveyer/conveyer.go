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
	size       int
	channels   map[string]chan string
	decorators []struct {
		fn            DecoratorFunc
		input, output string
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
	mu sync.Mutex
}

type Conveyer interface {
	RegisterDecorator(fn DecoratorFunc, input, output string)
	RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string)
	RegisterSeparator(fn SeparatorFunc, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

func New(size int) Conveyer {
	return &conveyerImpl{
		size:         size,
		channels:     make(map[string]chan string),
		decorators:   nil,
		multiplexers: nil,
		separators:   nil,
		mu:           sync.Mutex{},
	}
}

func (c *conveyerImpl) getChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, exists := c.channels[name]; exists {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) isChannelKnown(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, known := c.channels[name]
	return known
}

func (c *conveyerImpl) getAllChannelNames() map[string]bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	names := make(map[string]bool, len(c.channels))
	for name := range c.channels {
		names[name] = true
	}
	return names
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

	var runWaitGroup sync.WaitGroup
	errChannel := make(chan error, 1)

	for _, decorator := range c.decorators {
		runWaitGroup.Add(1)
		go func(currentDecorator struct {
			fn     DecoratorFunc
			input  string
			output string
		}) {
			defer runWaitGroup.Done()
			inputChannel := c.getChannel(currentDecorator.input)
			outputChannel := c.getChannel(currentDecorator.output)
			if err := currentDecorator.fn(ctx, inputChannel, outputChannel); err != nil {
				select {
				case errChannel <- err:
				default:
				}
			}
		}(decorator)
	}

	for _, multiplexer := range c.multiplexers {
		runWaitGroup.Add(1)
		go func(currentMultiplexer struct {
			fn     MultiplexerFunc
			inputs []string
			output string
		}) {
			defer runWaitGroup.Done()
			var inputChannels []chan string
			for _, inputName := range currentMultiplexer.inputs {
				channel := c.getChannel(inputName)
				inputChannels = append(inputChannels, channel)
			}
			outputChannel := c.getChannel(currentMultiplexer.output)
			if err := currentMultiplexer.fn(ctx, inputChannels, outputChannel); err != nil {
				select {
				case errChannel <- err:
				default:
				}
			}
		}(multiplexer)
	}

	for _, separator := range c.separators {
		runWaitGroup.Add(1)
		go func(currentSeparator struct {
			fn      SeparatorFunc
			input   string
			outputs []string
		}) {
			defer runWaitGroup.Done()
			inputChannel := c.getChannel(currentSeparator.input)
			var outputChannels []chan string
			for _, outputName := range currentSeparator.outputs {
				channel := c.getChannel(outputName)
				outputChannels = append(outputChannels, channel)
			}
			if err := currentSeparator.fn(ctx, inputChannel, outputChannels); err != nil {
				select {
				case errChannel <- err:
				default:
				}
			}
		}(separator)
	}

	done := make(chan struct{})
	go func() {
		runWaitGroup.Wait()
		close(done)
	}()

	select {
	case err := <-errChannel:
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
	if !c.isChannelKnown(input) {
		return ErrChanNotFound
	}
	ch := c.getChannel(input)
	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	if !c.isChannelKnown(output) {
		return "", ErrChanNotFound
	}
	ch := c.getChannel(output)
	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return data, nil
}
