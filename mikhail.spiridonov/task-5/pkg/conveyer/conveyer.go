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