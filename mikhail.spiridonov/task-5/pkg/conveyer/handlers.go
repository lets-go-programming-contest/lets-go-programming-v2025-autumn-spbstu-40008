package conveyer

import "context"

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