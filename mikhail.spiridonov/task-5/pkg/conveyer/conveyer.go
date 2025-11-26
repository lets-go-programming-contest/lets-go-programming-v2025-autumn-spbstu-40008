package conveyer

import (
    "context"
    "errors"
    "fmt"
    "sync"

    "golang.org/x/sync/errgroup"
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
    mu       sync.RWMutex
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
    c.mu.RLock()
    defer c.mu.RUnlock()

    ch, exists := c.channels[input]
    if !exists {

        return fmt.Errorf("%w: %s", ErrChanNotFound, input)
    }

    select {
    case ch <- data:

        return nil
    default:

        return errors.New("channel is full")
    }
}

func (c *DefaultConveyer) Recv(output string) (string, error) {
    c.mu.RLock()
    ch, exists := c.channels[output]
    c.mu.RUnlock()

    if !exists {

        return "", fmt.Errorf("%w: %s", ErrChanNotFound, output)
    }

    select {
    case data, ok := <-ch:
        if !ok {

            return undefined, nil
        }

        return data, nil
    default:

        return "", errors.New("no data available")
    }
}

func (c *DefaultConveyer) RegisterDecorator(
    fn func(ctx context.Context, input chan string, output chan string) error,
    input string,
    output string,
) {
    inCh := c.getOrCreateChannel(input)
    outCh := c.getOrCreateChannel(output)
    
    c.handlers = append(c.handlers, &decorator{
        fn:     fn,
        input:  inCh,
        output: outCh,
    })
}

func (c *DefaultConveyer) RegisterMultiplexer(
    fn func(ctx context.Context, inputs []chan string, output chan string) error,
    inputs []string,
    output string,
) {
    inChs := make([]chan string, len(inputs))
    for i, inputName := range inputs {
        inChs[i] = c.getOrCreateChannel(inputName)
    }

    outCh := c.getOrCreateChannel(output)
    
    c.handlers = append(c.handlers, &multiplexer{
        fn:     fn,
        input:  inChs,
        output: outCh,
    })
}

func (c *DefaultConveyer) RegisterSeparator(
    fn func(ctx context.Context, input chan string, outputs []chan string) error,
    input string,
    outputs []string,
) {
    inCh := c.getOrCreateChannel(input)
    
    outChs := make([]chan string, len(outputs))
    for i, outputName := range outputs {
        outChs[i] = c.getOrCreateChannel(outputName)
    }
    
    c.handlers = append(c.handlers, &separator{
        fn:     fn,
        input:  inCh,
        output: outChs,
    })
}

func (c *DefaultConveyer) getOrCreateChannel(name string) chan string {
    c.mu.Lock()
    defer c.mu.Unlock()

    if ch, exists := c.channels[name]; exists {

        return ch
    }

    ch := make(chan string, c.size)
    c.channels[name] = ch
	
    return ch
}

func (c *DefaultConveyer) close() {
    c.mu.Lock()
    defer c.mu.Unlock()

    for name, ch := range c.channels {
        close(ch)
        delete(c.channels, name)
    }
}
