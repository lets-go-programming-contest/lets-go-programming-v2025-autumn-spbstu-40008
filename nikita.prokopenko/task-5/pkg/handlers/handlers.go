package handlers

import (
	"context"
	"strings"
)

type DecorationError string

func (e DecorationError) Error() string {
	return string(e)
}

const (
	decoratorPrefix   = "decorated: "
	noDecoratorText   = "no decorator"
	noMultiplexerText = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, noDecoratorText) {
				return DecorationError("can't be decorated")
			}

			if !strings.HasPrefix(value, decoratorPrefix) {
				value = decoratorPrefix + value
			}

			select {
			case output <- value:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, out := range outputs {
			close(out)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil

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
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	done := make(chan struct{})
	errs := make(chan error, len(inputs))

	for _, in := range inputs {
		go func(ch chan string) {
			defer func() {
				done <- struct{}{}
			}()

			for {
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case val, ok := <-ch:
					if !ok {
						return
					}
					if strings.Contains(val, noMultiplexerText) {
						continue
					}
					select {
					case <-ctx.Done():
						errs <- ctx.Err()
						return
					case output <- val:
					}
				}
			}
		}(in)
	}

	for i := 0; i < len(inputs); i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errs:
			if err != nil {
				return err
			}
		case <-done:
		}
	}

	return nil
}