package conveyer

import (
	"context"
)

type decorator struct {
	function DecoratorFunc
	input    chan string
	output   chan string
}

func (dec *decorator) run(ctx context.Context) error {
	return dec.function(ctx, dec.input, dec.output)
}

type multiplexer struct {
	function MultiplexerFunc
	input    []chan string
	output   chan string
}

func (mux *multiplexer) run(ctx context.Context) error {
	return mux.function(ctx, mux.input, mux.output)
}

type separator struct {
	function SeparatorFunc
	input    chan string
	output   []chan string
}

func (sep *separator) run(ctx context.Context) error {
	return sep.function(ctx, sep.input, sep.output)
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
