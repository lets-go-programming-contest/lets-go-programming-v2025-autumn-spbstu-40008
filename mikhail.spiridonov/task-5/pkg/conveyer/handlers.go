package conveyer

import (
	"context"
)

type decorator struct {
	fn     DecoratorFunc
	input  chan string
	output chan string
}

func (d *decorator) run(ctx context.Context) error {
	return d.fn(ctx, d.input, d.output)
}

type multiplexer struct {
	fn     MultiplexerFunc
	input  []chan string
	output chan string
}

func (m *multiplexer) run(ctx context.Context) error {
	return m.fn(ctx, m.input, m.output)
}

type separator struct {
	fn     SeparatorFunc
	input  chan string
	output []chan string
}

func (s *separator) run(ctx context.Context) error {
	return s.fn(ctx, s.input, s.output)
}

type DecoratorFunc func(
    ctx context.Context,
    input chan string,
    output chan string,
) error

type MultiplexerFunc func(
    ctx context.Context,
    inputs []chan string,
    output chan string,
) error

type SeparatorFunc func(
    ctx context.Context,
    input chan string,
    outputs []chan string,
) error
