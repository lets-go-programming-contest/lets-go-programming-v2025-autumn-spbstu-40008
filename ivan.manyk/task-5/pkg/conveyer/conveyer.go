package conveyer

import (
	"context"
	"fmt"
	"sync"
)

const undefined = "undefined"

type conveyerType struct {
	size int
	channels map[string]chan string
	mu       sync.RWMutex
	handlers []func(ctx context.Context) error
}

func New(size int) *conveyerType {
    return &conveyerType{
        size:     size,
        channels: make(map[string]chan string),
    }
}

func (c *conveyerType) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyerType) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, exists := c.channels[name]
	return ch, exists
}

func (c *conveyerType) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputCh := c.getOrCreateChannel(input)
	outputCh := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputCh)
	})
}

func (c *conveyerType) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChs := make([]chan string, len(inputs))
	for i, input := range inputs {
		inputChs[i] = c.getOrCreateChannel(input)
	}
	outputCh := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChs, outputCh)
	})
}

func (c *conveyerType) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputCh := c.getOrCreateChannel(input)
	outputChs := make([]chan string, len(outputs))
	for i, output := range outputs {
		outputChs[i] = c.getOrCreateChannel(output)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputChs)
	})
}

func (c *conveyerType) Run(ctx context.Context) error {
    var wg sync.WaitGroup
    errChan := make(chan error, 1)
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    for _, handler := range c.handlers {
        wg.Add(1)
        go func(h func(ctx context.Context) error) {
            defer wg.Done()
            if err := h(ctx); err != nil {
                select {
                case errChan <- err:
                default:
                }
                cancel()
            }
        }(handler)
    }

    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-ctx.Done():
        wg.Wait()
        c.closeAllChannels()
        return nil
    case err := <-errChan:
        wg.Wait()
        c.closeAllChannels()
        return err
    case <-done:
        c.closeAllChannels()
        return nil
    }
}

func (c *conveyerType) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *conveyerType) Send(input string, data string) error {
	ch, exists := c.getChannel(input)
	if !exists {
		return fmt.Errorf("chan not found")
	}

	ch <- data
	return nil
}

func (c *conveyerType) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)
	if !exists {
		return "", fmt.Errorf("chan not found")
	}

	data, ok := <-ch
	if !ok {
		return undefined, nil
	}
	return data, nil
}
