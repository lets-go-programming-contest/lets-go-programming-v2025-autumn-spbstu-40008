package handlers

import (
	"context"
	"strings"
)

type DecoratorError string

func (e DecoratorError) Error() string {
	return string(e)
}

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	defer close(output)

	for {
		select {
		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return DecoratorError("can't be decorated")
			}

			if !strings.HasPrefix(value, "decorated: ") {
				value = "decorated: " + value
			}

			select {
			case output <- value:
			case <-ctx.Done():
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	for _, out := range outputs {
		defer close(out)
	}

	if len(outputs) == 0 {
		return nil
	}

	idx := 0

	for {
		select {
		case value, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[idx] <- value:
				idx = (idx + 1) % len(outputs)
			case <-ctx.Done():
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	type result struct {
		value string
		err   error
	}

	results := make(chan result, len(inputs))

	for _, in := range inputs {
		go func(ch chan string) {
			for {
				select {
				case val, ok := <-ch:
					if !ok {
						results <- result{err: nil}
						return
					}
					if strings.Contains(val, "no multiplexer") {
						continue
					}
					select {
					case output <- val:
					case <-ctx.Done():
						results <- result{err: ctx.Err()}
						return
					}
				case <-ctx.Done():
					results <- result{err: ctx.Err()}
					return
				}
			}
		}(in)
	}

	for i := 0; i < len(inputs); i++ {
		res := <-results
		if res.err != nil {
			return res.err
		}
	}

	return nil
}