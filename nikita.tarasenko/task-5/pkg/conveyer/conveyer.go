package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Интерфейс из задания.
type conveyer interface {
	RegisterDecorator(fn func(context.Context, chan string, chan string) error, input, output string)
	RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(context.Context, chan string, []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

var ErrChannelNotFound = errors.New("channel not found")

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
	}
}

func (p *pipeline) getOrCreateChannel(name string) chan string {
	p.mu.RLock()
	if ch, ok := p.channels[name]; ok {
		p.mu.RUnlock()
		return ch
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()
	if ch, ok := p.channels[name]; ok {
		return ch
	}
	ch := make(chan string, p.size)
	p.channels[name] = ch
	return ch
}

func (p *pipeline) getChannel(name string) (chan string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	ch, ok := p.channels[name]
	return ch, ok
}

func (p *pipeline) closeAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, ch := range p.channels {
		close(ch)
	}
}

func (p *pipeline) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input, output string,
) {
	in := p.getOrCreateChannel(input)
	out := p.getOrCreateChannel(output)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (p *pipeline) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	ins := make([]chan string, len(inputs))
	for i, name := range inputs {
		ins[i] = p.getOrCreateChannel(name)
	}
	out := p.getOrCreateChannel(output)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (p *pipeline) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	outs := make([]chan string, len(outputs))
	for i, name := range outputs {
		outs[i] = p.getOrCreateChannel(name)
	}
	in := p.getOrCreateChannel(input)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outs)
	})
}

func (p *pipeline) Run(parentCtx context.Context) error {
	defer p.closer.Do(p.closeAll)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	p.mu.RLock()
	for _, h := range p.handlers {
		h := h
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				cancel()
			}
		}()
	}
	p.mu.RUnlock()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("pipeline failed: %w", err)
	case <-done:
		return nil
	case <-ctx.Done():
		select {
		case err := <-errCh:
			return fmt.Errorf("pipeline failed: %w", err)
		default:
			return nil
		}
	}
}

func (p *pipeline) Send(chName string, data string) error {
	ch, ok := p.getChannel(chName)
	if !ok {
		return fmt.Errorf("channel %q not found: %w", chName, ErrChannelNotFound)
	}
	ch <- data
	return nil
}

func (p *pipeline) Recv(chName string) (string, error) {
	ch, ok := p.getChannel(chName)
	if !ok {
		return "", fmt.Errorf("channel %q not found: %w", chName, ErrChannelNotFound)
	}
	val, ok := <-ch
	if !ok {
		return undefinedValue, nil
	}
	return val, nil
}
