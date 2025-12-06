package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	noDecoratorMessage   = "no decorator"
	decoratorPrefix      = "decorated: "
	noMultiplexerMessage = "no multiplexer"
)

func drainInput(ctx context.Context, input chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case _, ok := <-input:
			if !ok {
				return nil
			}
		}
	}
}

func runPrefixDecorator(ctx context.Context, input, output chan string) error {
	defer func() {
		if output != nil {
			close(output)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, noDecoratorMessage) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(value, decoratorPrefix) {
				value = decoratorPrefix + value
			}

			if output != nil {
				select {
				case output <- value:
				case <-ctx.Done():
					return nil
				}
			}
		}
	}
}

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	if output == nil {
		return drainInput(ctx, input)
	}
	return runPrefixDecorator(ctx, input, output)
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, outCh := range outputs {
			if outCh != nil {
				select {
				case <-outCh:
				default:
					close(outCh)
				}
			}
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

			outCh := outputs[idx]
			if outCh != nil {
				select {
				case outCh <- value:
					idx = (idx + 1) % len(outputs)
				case <-ctx.Done():
					return nil
				}
			} else {
				idx = (idx + 1) % len(outputs)
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if output == nil {
		return nil
	}

	defer func() {
		if output != nil {
			select {
			case <-output:
			default:
				close(output)
			}
		}
	}()

	if len(inputs) == 0 {
		return nil
	}

	done := make(chan struct{})
	defer close(done)

	for _, inputCh := range inputs {
		go func(ch chan string) {
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case val, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(val, noMultiplexerMessage) {
						continue
					}

					if output != nil {
						select {
						case output <- val:
						case <-ctx.Done():
							return
						case <-done:
							return
						}
					}
				}
			}
		}(inputCh)
	}

	<-ctx.Done()
	return nil
}