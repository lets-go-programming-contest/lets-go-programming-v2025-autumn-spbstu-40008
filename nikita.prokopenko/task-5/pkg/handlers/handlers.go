package handlers

import (
	"context"
	"strings"
	"sync"
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

		case value, open := <-input:
			if !open {
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

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, open := <-input:
			if !open {
				return nil
			}

			select {
			case outputs[counter] <- value:
				counter = (counter + 1) % len(outputs)
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

	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(inputChan chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, open := <-inputChan:
					if !open {
						return
					}

					if strings.Contains(value, noMultiplexerText) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- value:
					}
				}
			}
		}(in)
	}

	wg.Wait()
	return nil
}