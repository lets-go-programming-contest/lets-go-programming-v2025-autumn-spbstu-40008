package conveyer

import (
	"context"
	"errors"
	"time"
)

type Conveyer[T any] struct {
	input  chan T
	output chan T
}

func NewConveyer[T any](size int) *Conveyer[T] {
	return &Conveyer[T]{
		input:  make(chan T, size),
		output: make(chan T, size),
	}
}

func (c *Conveyer[T]) Send(ctx context.Context, data T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.input <- data:
		return nil
	case <-time.After(200 * time.Millisecond):
		return errors.New("send timeout")
	}
}

func (c *Conveyer[T]) Recv(ctx context.Context) (T, error) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case v := <-c.output:
		return v, nil
	case <-time.After(200 * time.Millisecond):
		return zero, errors.New("recv timeout")
	}
}

func (c *Conveyer[T]) Input() chan T {
	return c.input
}

func (c *Conveyer[T]) Output() chan T {
	return c.output
}
