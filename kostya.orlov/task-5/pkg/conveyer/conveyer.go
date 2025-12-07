package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type ConveyerFunc func(ctx context.Context, input chan string, outputs []chan string) error
type MultiplexerFuncAdaptor func(ctx context.Context, inputs []chan string, outputs []chan string) error

type HandlerRegistration struct {
	Type        string
	Fn          ConveyerFunc
	InputID     string
	InputIDs    []string
	OutputIDs   []string
	InputChans  []chan string
	OutputChans []chan string
}

type Conveyer struct {
	channels map[string]chan string
	handlers []HandlerRegistration

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup
	mu sync.RWMutex

	errorChan chan error
	size      int
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels:  make(map[string]chan string),
		handlers:  make([]HandlerRegistration, 0),
		errorChan: make(chan error, 1),
		size:      size,
	}
}

func (c *Conveyer) AddChannel(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.channels[id]; exists {
		return fmt.Errorf("channel with ID %s already exists", id)
	}
	c.channels[id] = make(chan string, c.size)
	return nil
}

func (c *Conveyer) RegisterDecorator(fn ConveyerFunc, inputID string, outputID string) {
	c.handlers = append(c.handlers, HandlerRegistration{
		Type:      "Decorator",
		Fn:        fn,
		InputID:   inputID,
		OutputIDs: []string{outputID},
	})
}

func (c *Conveyer) RegisterMultiplexer(fn ConveyerFunc, inputIDs []string, outputID string) {
	c.handlers = append(c.handlers, HandlerRegistration{
		Type:      "Multiplexer",
		Fn:        fn,
		InputIDs:  inputIDs,
		OutputIDs: []string{outputID},
	})
}

func (c *Conveyer) RegisterSeparator(fn ConveyerFunc, inputID string, outputIDs []string) {
	c.handlers = append(c.handlers, HandlerRegistration{
		Type:      "Separator",
		Fn:        fn,
		InputID:   inputID,
		OutputIDs: outputIDs,
	})
}

func (c *Conveyer) resolveChannels() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := range c.handlers {
		handler := &c.handlers[i]

		if handler.Type == "Multiplexer" {
			handler.InputChans = make([]chan string, 0, len(handler.InputIDs))
			for _, inputID := range handler.InputIDs {
				if ch, ok := c.channels[inputID]; ok {
					handler.InputChans = append(handler.InputChans, ch)
				} else {
					return fmt.Errorf("input channel ID %s not found for Multiplexer", inputID)
				}
			}
		} else if handler.InputID != "" {
			if ch, ok := c.channels[handler.InputID]; ok {
				handler.InputChans = []chan string{ch}
			} else {
				return fmt.Errorf("input channel ID %s not found for handler %s", handler.InputID, handler.Type)
			}
		}

		handler.OutputChans = make([]chan string, 0, len(handler.OutputIDs))
		for _, outputID := range handler.OutputIDs {
			if ch, ok := c.channels[outputID]; ok {
				handler.OutputChans = append(handler.OutputChans, ch)
			} else {
				return fmt.Errorf("output channel ID %s not found for handler %s", outputID, handler.Type)
			}
		}
	}
	return nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()

	if err := c.resolveChannels(); err != nil {
		return fmt.Errorf("failed to resolve channels: %w", err)
	}

	for _, handler := range c.handlers {
		c.wg.Add(1)
		go func(h HandlerRegistration) {
			defer c.wg.Done()

			var inputChan chan string
			if len(h.InputChans) > 0 {
				inputChan = h.InputChans[0]
			}

			err := h.Fn(c.ctx, inputChan, h.OutputChans)

			if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
				select {
				case c.errorChan <- fmt.Errorf("handler %s failed: %w", h.Type, err):
				case <-c.ctx.Done():
				}
			}
		}(handler)
	}

	var runError error

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case err := <-c.errorChan:
		c.cancel()
		runError = err

	case <-done:
		runError = nil

	case <-ctx.Done():
		c.cancel()
		runError = ctx.Err()
	}

	<-done

	c.closeAllChannels()

	return runError
}

func (c *Conveyer) Send(inputID string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[inputID]
	c.mu.RUnlock()

	if !ok {
		return errors.New("chan not found")
	}

	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case ch <- data:
		return nil
	}
}

func (c *Conveyer) Recv(outputID string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[outputID]
	c.mu.RUnlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	select {
	case <-c.ctx.Done():
		return "", c.ctx.Err()

	case data, open := <-ch:
		if !open {
			return "undefined", nil
		}
		return data, nil
	}
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, ch := range c.channels {
		select {
		case <-c.ctx.Done():
			return
		default:
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				close(ch)
			}()
		}
	}
}
