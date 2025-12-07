package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
)

type conveyer interface {
	RegisterDecorator(
		handlerFunc func(context.Context, chan string, chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		handlerFunc func(context.Context, []chan string, chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		handlerFunc func(context.Context, chan string, []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

const undefinedValue = "undefined"

type pipeline struct {
	size     int
	mu       sync.RWMutex
	channels map[string]chan string
	handlers []func(context.Context) error
	closer   sync.Once
}

func New(size int) conveyer {
	return &pipeline{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
		mu:       sync.RWMutex{},
		closer:   sync.Once{},
	}
}

func (p *pipeline) getOrCreateChannel(name string) chan string {
	p.mu.RLock()

	channel, ok := p.channels[name]
	if ok {
		p.mu.RUnlock()
		return channel
	}

	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	channel, ok = p.channels[name]
	if ok {
		return channel
	}

	channel = make(chan string, p.size)
	p.channels[name] = channel

	return channel
}

func (p *pipeline) getChannel(name string) (chan string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	channel, ok := p.channels[name]

	return channel, ok
}

func (p *pipeline) closeAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, channel := range p.channels {
		close(channel)
	}
}

func (p *pipeline) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inputChannel := p.getOrCreateChannel(input)
	outputChannel := p.getOrCreateChannel(output)

	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannel)
	})
}

func (p *pipeline) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))

	for i, name := range inputs {
		inputChannels[i] = p.getOrCreateChannel(name)
	}

	outputChannel := p.getOrCreateChannel(output)

	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannels, outputChannel)
	})
}

func (p *pipeline) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	outputChannels := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputChannels[i] = p.getOrCreateChannel(name)
	}

	inputChannel := p.getOrCreateChannel(input)

	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannels)
	})
}

func (p *pipeline) Run(parentCtx context.Context) error {
	defer p.closer.Do(p.closeAll)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	var waitGroup sync.WaitGroup
	errorChannel := make(chan error, 1)

	p.mu.RLock()

	for _, handlerFunc := range p.handlers {
		waitGroup.Add(1)

		go func(handler func(context.Context) error) {
			defer waitGroup.Done()

			if err := handler(ctx); err != nil {
				select {
				case errorChannel <- err:
				default:
				}

				cancel()
			}
		}(handlerFunc)
	}

	p.mu.RUnlock()

	done := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case err := <-errorChannel:
		return fmt.Errorf("pipeline failed: %w", err)

	case <-done:
		return nil

	case <-ctx.Done():
		select {
		case err := <-errorChannel:
			return fmt.Errorf("pipeline failed: %w", err)
		default:
			return nil
		}
	}
}

func (p *pipeline) Send(channelName string, data string) error {
	channel, ok := p.getChannel(channelName)
	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (p *pipeline) Recv(channelName string) (string, error) {
	channel, ok := p.getChannel(channelName)
	if !ok {
		return "", ErrChanNotFound
	}

	value, ok := <-channel
	if !ok {
		return undefinedValue, nil
	}

	return value, nil
}
