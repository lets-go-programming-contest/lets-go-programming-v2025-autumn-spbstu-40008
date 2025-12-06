// pkg/handlers/handlers.go
package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	prefix := "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(value, prefix) {
				value = prefix + value
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
		for _, ch := range outputs {
			close(ch)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	idx := 0
	total := len(outputs)

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
				idx = (idx + 1) % total
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if output == nil {
		return nil
	}

	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	processChannel := func(ch chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-ch:
				if !ok {
					return
				}

				if strings.Contains(value, "no multiplexer") {
					continue
				}

				select {
				case output <- value:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, ch := range inputs {
		wg.Add(1)
		go processChannel(ch)
	}

	wg.Wait()
	return nil
}