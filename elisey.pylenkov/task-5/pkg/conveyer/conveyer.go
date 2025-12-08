package conveyer

import (
	"context"
	"errors"
	"sync"
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
	wg       sync.WaitGroup
	running  bool
	ready    chan struct{}
	errors   chan error
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		ready:    make(chan struct{}),
		errors:   make(chan error, 1),
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input, output string,
) {
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
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string, output string,
) {
	for _, input := range inputs {
		c.getOrCreateChannel(input)
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
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string, outputs []string,
) {
	c.getOrCreateChannel(input)
	for _, output := range outputs {
		c.getOrCreateChannel(output)
	}

	c.handlers = append(c.handlers, handlerFunc{
		fn:      fn,
		inputs:  []string{input},
		outputs: outputs,
		htype:   separatorType,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return errors.New("conveyer already running")
	}
	c.running = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
	}()

	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	defer cancel()

	// Собираем все каналы
	c.mu.RLock()
	channels := make(map[string]chan string)
	for name, ch := range c.channels {
		channels[name] = ch
	}
	c.mu.RUnlock()

	c.wg.Add(len(c.handlers))

	// Запускаем обработчики
	for _, h := range c.handlers {
		go func(h handlerFunc) {
			defer c.wg.Done()

			// Получаем каналы
			inputs := make([]chan string, len(h.inputs))
			for i, name := range h.inputs {
				inputs[i] = channels[name]
			}

			outputs := make([]chan string, len(h.outputs))
			for i, name := range h.outputs {
				outputs[i] = channels[name]
			}

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
				case c.errors <- err:
				default:
				}
			}
		}(h)
	}

	// Сигнализируем, что обработчики запущены
	close(c.ready)

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.closeAllChannels()
		c.wg.Wait()
		return ctx.Err()
	case err := <-c.errors:
		cancel()
		c.closeAllChannels()
		c.wg.Wait()
		return err
	case <-done:
		c.closeAllChannels()
		return nil
	}
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return errors.New("chan not found")
	}

	// Отправка с возможностью блокировки до отправки
	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	// Чтение с блокировкой до получения данных или закрытия канала
	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return data, nil
}

func (c *Conveyer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
}
