package conveyer

import (
    "context"
    "errors"
    "fmt"
)

const (
    undefined = "undefined"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer interface {
    RegisterDecorator(
        fn func(ctx context.Context, input chan string, output chan string) error,
        input string,
        output string,
    )
    RegisterMultiplexer(
        fn func(ctx context.Context, inputs []chan string, output chan string) error,
        inputs []string,
        output string,
    )
    RegisterSeparator(
        fn func(ctx context.Context, input chan string, outputs []chan string) error,
        input string,
        outputs []string,
    )
    Run(ctx context.Context) error
    Send(input string, data string) error
    Recv(output string) (string, error)
}

type DefaultConveyer struct {
    size     int
    channels map[string]chan string
    handlers []handler
}

type handler interface {
    run(ctx context.Context) error
}

func New(size int) *DefaultConveyer {
    return &DefaultConveyer{
        size:     size,
        channels: make(map[string]chan string),
        handlers: make([]handler, 0),
    }
}

func (c *DefaultConveyer) Run(ctx context.Context) error {
    defer c.close()

    g, gCtx := errgroup.WithContext(ctx)

    for _, h := range c.handlers {
        handler := h
        g.Go(func() error {
            return handler.run(gCtx)
        })
    }

    if err := g.Wait(); err != nil {
        return fmt.Errorf("conveyer stopped: %w", err)
    }

    return nil
}

func (c *DefaultConveyer) Send(input string, data string) error {
    ch, exists := c.channels[input]
    if !exists {
        return fmt.Errorf("%w: %s", ErrChanNotFound, input)
    }

    ch <- data
    return nil
}

func (c *DefaultConveyer) Recv(output string) (string, error) {
    ch, exists := c.channels[output]
    if !exists {
        return "", fmt.Errorf("%w: %s", ErrChanNotFound, output)
    }

    data, ok := <-ch
    if !ok {
        return undefined, nil
    }

    return data, nil
}

func (c *DefaultConveyer) getOrCreateChannel(name string) chan string {
    if ch, exists := c.channels[name]; exists {
        return ch
    }

    ch := make(chan string, c.size)
    c.channels[name] = ch
    return ch
}

func (c *DefaultConveyer) close() {
    for _, ch := range c.channels {
        close(ch)
    }
}