package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

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
				return fmt.Errorf("%w: %s", ErrCantBeDecorated, data)
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
			hasActiveChannels := false
			dataSent := false

			for index := range inputs {
				select {
				case <-ctx.Done():
					return fmt.Errorf("context done: %w", ctx.Err())

				default:
					select {
					case data, ok := <-inputs[index]:
						if !ok {
							continue
						}

						hasActiveChannels = true 

						if strings.Contains(data, skipSubstring) {
							continue
						}

						select {
						case output <- data:
						case <-ctx.Done():
							return fmt.Errorf("context done: %w", ctx.Err())
						}

					default:
						hasActiveChannels = true
					}
				}
			}

			if !hasActiveChannels {
				return nil
			}

			if !dataSent {
				select {
				case <-ctx.Done():
					return fmt.Errorf("context done: %w", ctx.Err())

				default:
				}
			}
		}
	}		
}