package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound  = errors.New("chan not found")
	ErrChannelExists = errors.New("channel already exists")
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type HandlerRegistration struct {
	Type        string
	Fn          interface{}
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
		errorChan: make(chan error, 10),
		size:      size,
		waitGroup: sync.WaitGroup{},
		mu:        sync.RWMutex{},
	}
}

func (c *Conveyer) AddChannel(channelID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.channels[channelID]; exists {
		return ErrChannelExists
	}

	c.channels[channelID] = make(chan string, c.size)
	return nil
}

func (c *Conveyer) RegisterDecorator(fn interface{}, inputID string, outputID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.channels[inputID]; !exists {
		c.channels[inputID] = make(chan string, c.size)
	}
	if _, exists := c.channels[outputID]; !exists {
		c.channels[outputID] = make(chan string, c.size)
	}

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

func (c *Conveyer) RegisterMultiplexer(fn interface{}, inputIDs []string, outputID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, inputID := range inputIDs {
		if _, exists := c.channels[inputID]; !exists {
			c.channels[inputID] = make(chan string, c.size)
		}
	}
	if _, exists := c.channels[outputID]; !exists {
		c.channels[outputID] = make(chan string, c.size)
	}

	c.handlers = append(c.handlers, HandlerRegistration{
		Type:        "Multiplexer",
		Fn:          fn,
		InputIDs:    inputIDs,
		OutputIDs:   []string{outputID},
		InputID:     "",
		InputChans:  nil,
		OutputChans: nil,
	})
}

func (c *Conveyer) RegisterSeparator(fn interface{}, inputID string, outputIDs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.channels[inputID]; !exists {
		c.channels[inputID] = make(chan string, c.size)
	}
	for _, outputID := range outputIDs {
		if _, exists := c.channels[outputID]; !exists {
			c.channels[outputID] = make(chan string, c.size)
		}
	}

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
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.handlers {
		handler := &c.handlers[i]

		if handler.Type == "Multiplexer" {
			handler.InputChans = make([]chan string, 0, len(handler.InputIDs))
			for _, inputID := range handler.InputIDs {
				handler.InputChans = append(handler.InputChans, c.channels[inputID])
			}
		} else if handler.InputID != "" {
			handler.InputChans = []chan string{c.channels[handler.InputID]}
		}

		handler.OutputChans = make([]chan string, 0, len(handler.OutputIDs))
		for _, outputID := range handler.OutputIDs {
			handler.OutputChans = append(handler.OutputChans, c.channels[outputID])
		}
	}
	return nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)

	if err := c.resolveChannels(); err != nil {
		c.cancel()
		return err
	}

	for _, handler := range c.handlers {
		c.waitGroup.Add(1)
		go func(handlerReg HandlerRegistration) {
			defer c.waitGroup.Done()

			var err error
			if handlerReg.Type == "Multiplexer" {
				if len(handlerReg.OutputChans) > 0 {
					if fn, ok := handlerReg.Fn.(MultiplexerFunc); ok {
						err = fn(c.ctx, handlerReg.InputChans, handlerReg.OutputChans[0])
					} else if fn, ok := handlerReg.Fn.(func(context.Context, []chan string, chan string) error); ok {
						err = fn(c.ctx, handlerReg.InputChans, handlerReg.OutputChans[0])
					} else {
						err = fmt.Errorf("unsupported multiplexer function type")
					}
				} else {
					err = errors.New("multiplexer requires output channel")
				}
			} else if handlerReg.Type == "Decorator" {
				var inputChan, outputChan chan string
				if len(handlerReg.InputChans) > 0 {
					inputChan = handlerReg.InputChans[0]
				}
				if len(handlerReg.OutputChans) > 0 {
					outputChan = handlerReg.OutputChans[0]
				}

				if fn, ok := handlerReg.Fn.(DecoratorFunc); ok {
					err = fn(c.ctx, inputChan, outputChan)
				} else if fn, ok := handlerReg.Fn.(func(context.Context, chan string, chan string) error); ok {
					err = fn(c.ctx, inputChan, outputChan)
				} else {
					err = fmt.Errorf("unsupported decorator function type")
				}
			} else if handlerReg.Type == "Separator" {
				var inputChan chan string
				if len(handlerReg.InputChans) > 0 {
					inputChan = handlerReg.InputChans[0]
				}

				if fn, ok := handlerReg.Fn.(SeparatorFunc); ok {
					err = fn(c.ctx, inputChan, handlerReg.OutputChans)
				} else if fn, ok := handlerReg.Fn.(func(context.Context, chan string, []chan string) error); ok {
					err = fn(c.ctx, inputChan, handlerReg.OutputChans)
				} else {
					err = fmt.Errorf("unsupported separator function type")
				}
			} else {
				err = fmt.Errorf("unknown handler type: %s", handlerReg.Type)
			}

			if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				select {
				case c.errorChan <- err:
				default:
				}
				c.cancel()
			}
		}(handler)
	}

	go func() {
		c.waitGroup.Wait()
		close(c.errorChan)
	}()

	var runError error

	select {
	case err, ok := <-c.errorChan:
		if ok {
			runError = err
		}
	case <-ctx.Done():
		runError = ctx.Err()
		c.cancel()
	}


	if runError != nil {
		for range c.errorChan {} 
	}
    
	c.cancel() 

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

	if c.ctx == nil {
		channel <- data
		return nil
	}

	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
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

	if c.ctx == nil {
		data, open := <-channel
		if !open {
			return "undefined", nil
		}
		return data, nil
	}

	select {
	case <-c.ctx.Done():
		return "", c.ctx.Err()
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

	closedChannels := make(map[chan string]bool)

	for _, channel := range c.channels {
		if closedChannels[channel] {
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {}
			}()
			close(channel)
			closedChannels[channel] = true
		}()
	}
}