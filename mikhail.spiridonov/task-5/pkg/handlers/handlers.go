package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {

	const prefix = "decorated: "

	const errorSubstring = "no decorator"

	for {
		select {
		case <-ctx.Done():

			return fmt.Errorf("context done: %w", fmt.Errorf("context done: %w", ctx.Err()))
		case data, ok := <-input:
			if !ok {

				return nil
			}

			if strings.Contains(data, errorSubstring) {

				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return fmt.Errorf("context done: %w", ctx.Err())
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())

		case data, ok := <-input:
			if !ok {
				for _, out := range outputs {
					close(out)
				}

				return nil
			}

			idx := counter % len(outputs)
			counter++

			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return fmt.Errorf("context done: %w", ctx.Err())
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {

	const skipSubstring = "no multiplexer"

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())

		default:
			processed := false

			for _, inputCh := range inputs {
				select {
				case data, ok := <-inputCh:
					if !ok {
						continue
					}
					processed = true

					if strings.Contains(data, skipSubstring) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return fmt.Errorf("context done: %w", ctx.Err())
					}

				default:

				}
			}

			if !processed {
				select {
				case <-ctx.Done():
					return fmt.Errorf("context done: %w", ctx.Err())
				}
			}
		}
	}
}
