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

    )
    RegisterMultiplexer(

    )
    RegisterSeparator(

    )
    Run(ctx context.Context) error
    Send(input string, data string) error
    Recv(output string) (string, error)
}

type DefaultConveyer struct {

}

type handler interface {

}

func New(size int) *DefaultConveyer {
    return &DefaultConveyer{

    }
}