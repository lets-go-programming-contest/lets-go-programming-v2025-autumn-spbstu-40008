package handlers

import "context"

type DecoratorFunc[T any] func(context.Context, T) (T, error)

type Decorator[T any] struct {
	fn     DecoratorFunc[T]
	input  chan T
	output chan T
}

func NewDecorator[T any](fn DecoratorFunc[T], size int) *Decorator[T] {
	return &Decorator[T]{
		fn:     fn,
		input:  make(chan T, size),
		output: make(chan T, size),
	}
}

func (d *Decorator[T]) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(d.output)
			return
		case v, ok := <-d.input:
			if !ok {
				close(d.output)
				return
			}
			res, err := d.fn(ctx, v)
			if err == nil {
				select {
				case d.output <- res:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func (d *Decorator[T]) Input() chan T {
	return d.input
}

func (d *Decorator[T]) Output() chan T {
	return d.output
}
