package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound  = errors.New("channel not found")
	ErrChannelExists = errors.New("channel already exists")
)

type ConveyerFunc func(ctx context.Context, input chan string, outputs []chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, outputs []chan string) error

type HandlerRegistration struct {
	Type        string
	Fn          ConveyerFunc
	MultiplexFn MultiplexerFunc
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

	waitGroup sync.WaitGroup
	mu        sync.RWMutex

	errorChan chan error
	size      int
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels:  make(map[string]chan string),
		handlers:  make([]HandlerRegistration, 0),
		errorChan: make(chan error, 1),
		size:      size,
		waitGroup: sync.WaitGroup{},
		mu:        sync.RWMutex{},
	}
}

func (c *Conveyer) AddChannel(channelID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.channels[channelID]; exists {
		return fmt.Errorf("%w: ID %s", ErrChannelExists, channelID)
	}

	c.channels[channelID] = make(chan string, c.size)
	return nil
}

func (c *Conveyer) RegisterDecorator(fn ConveyerFunc, inputID string, outputID string) {
	c.handlers = append(c.handlers, HandlerRegistration{
		Type:        "Decorator",
		Fn:          fn,
		InputID:     inputID,
		OutputIDs:   []string{outputID},
		InputIDs:    nil,
		InputChans:  nil,
		OutputChans: nil,
	})
}

func (c *Conveyer) RegisterMultiplexer(fn MultiplexerFunc, inputIDs []string, outputID string) {
	c.handlers = append(c.handlers, HandlerRegistration{
		Type:        "Multiplexer",
		MultiplexFn: fn,
		InputIDs:    inputIDs,
		OutputIDs:   []string{outputID},
		InputID:     "",
		InputChans:  nil,
		OutputChans: nil,
	})
}

func (c *Conveyer) RegisterSeparator(fn ConveyerFunc, inputID string, outputIDs []string) {
	c.handlers = append(c.handlers, HandlerRegistration{
		Type:        "Separator",
		Fn:          fn,
		InputID:     inputID,
		OutputIDs:   outputIDs,
		InputIDs:    nil,
		InputChans:  nil,
		OutputChans: nil,
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
				if channel, ok := c.channels[inputID]; ok {
					handler.InputChans = append(handler.InputChans, channel)
				} else {
					return fmt.Errorf("input channel ID %s not found: %w", inputID, ErrChanNotFound)
				}
			}
		} else if handler.InputID != "" {
			if channel, ok := c.channels[handler.InputID]; ok {
				handler.InputChans = []chan string{channel}
			} else {
				return fmt.Errorf("input channel ID %s not found for handler %s: %w", handler.InputID, handler.Type, ErrChanNotFound)
			}
		}

		handler.OutputChans = make([]chan string, 0, len(handler.OutputIDs))
		for _, outputID := range handler.OutputIDs {
			if channel, ok := c.channels[outputID]; ok {
				handler.OutputChans = append(handler.OutputChans, channel)
			} else {
				return fmt.Errorf("output channel ID %s not found for handler %s: %w", outputID, handler.Type, ErrChanNotFound)
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
		c.waitGroup.Add(1)
		go func(handlerReg HandlerRegistration) {
			defer c.waitGroup.Done()

			var err error
			if handlerReg.Type == "Multiplexer" {
				err = handlerReg.MultiplexFn(c.ctx, handlerReg.InputChans, handlerReg.OutputChans)
			} else {
				var inputChannel chan string
				if len(handlerReg.InputChans) > 0 {
					inputChannel = handlerReg.InputChans[0]
				}
				err = handlerReg.Fn(c.ctx, inputChannel, handlerReg.OutputChans)
			}

			if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				select {
				case c.errorChan <- fmt.Errorf("handler %s failed: %w", handlerReg.Type, err):
				case <-c.ctx.Done():
				}
			}
		}(handler)
	}

	var runError error

	done := make(chan struct{})
	go func() {
		c.waitGroup.Wait()
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
		runError = fmt.Errorf("external context cancelled: %w", ctx.Err())
	}

	<-done

	c.closeAllChannels()

	return runError
}

func (c *Conveyer) Send(inputID string, data string) error {
	c.mu.RLock()
	channel, ok := c.channels[inputID]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	select {
	case <-c.ctx.Done():
		return fmt.Errorf("send failed: %w", c.ctx.Err())
	case channel <- data:
		return nil
	}
}

func (c *Conveyer) Recv(outputID string) (string, error) {
	c.mu.RLock()
	channel, ok := c.channels[outputID]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	select {
	case <-c.ctx.Done():
		return "", fmt.Errorf("receive failed: %w", c.ctx.Err())

	case data, open := <-channel:
		if !open {
			return "undefined", nil
		}
		return data, nil
	}
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, channel := range c.channels {
		select {
		case <-c.ctx.Done():
			return
		default:
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				close(channel)
			}()
		}
	}
}