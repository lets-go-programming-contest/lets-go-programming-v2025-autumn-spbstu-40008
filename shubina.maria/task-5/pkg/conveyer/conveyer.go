package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("chan not found")

// Conveyer — интерфейс, описанный в задании.
type Conveyer interface {
	RegisterDecorator(fn func(ctx context.Context, input, output chan string) error, input, output string)
	RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	chMutex  sync.RWMutex
	handlers []func(ctx context.Context) error
}

// New создаёт новый конвейер с указанным размером буферов каналов.
func New(size int) Conveyer {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
	}
}

// ensureChannel создаёт канал, если его ещё нет.
func (c *conveyerImpl) ensureChannel(name string) chan string {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()
	if ch, exists := c.channels[name]; exists {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

// getChannel возвращает существующий канал или nil.
func (c *conveyerImpl) getChannel(name string) chan string {
	c.chMutex.RLock()
	defer c.chMutex.RUnlock()
	return c.channels[name]
}

// RegisterDecorator регистрирует обработчик-модификатор.
func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inCh := c.ensureChannel(input)
	outCh := c.ensureChannel(output)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

// RegisterMultiplexer регистрирует мультиплексор.
func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputChans[i] = c.ensureChannel(name)
	}
	outCh := c.ensureChannel(output)
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChans, outCh)
	})
}

// RegisterSeparator регистрирует сепаратор.
func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.ensureChannel(input)
	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = c.ensureChannel(name)
	}
	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outChans)
	})
}

// Send отправляет данные в указанный канал.
func (c *conveyerImpl) Send(input, data string) error {
	ch := c.getChannel(input)
	if ch == nil {
		return ErrChannelNotFound
	}
	ch <- data
	return nil
}

// Recv получает данные из указанного канала.
func (c *conveyerImpl) Recv(output string) (string, error) {
	ch := c.getChannel(output)
	if ch == nil {
		return "", ErrChannelNotFound
	}
	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return data, nil
}

// Run запускает все зарегистрированные обработчики в отдельных горутинах.
func (c *conveyerImpl) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, h := range c.handlers {
		h := h // захват в замыкании
		g.Go(func() error { return h(ctx) })
	}
	err := g.Wait()
	c.closeAllChannels()
	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}
	return nil
}

// closeAllChannels закрывает все внутренние каналы.
func (c *conveyerImpl) closeAllChannels() {
	c.chMutex.Lock()
	defer c.chMutex.Unlock()
	for _, ch := range c.channels {
		close(ch)
	}
}
