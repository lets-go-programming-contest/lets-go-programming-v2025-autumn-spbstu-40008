// pkg/handlers/handlers.go
package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotDecorate = errors.New("can't be decorated")

const (
	decoratorPrefix      = "decorated: "
	decoratorSkipMessage = "no decorator"
	multiplexerSkip      = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	if output == nil {
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, isOpen := <-input:
				if !isOpen {
					return nil
				}
			}
		}
	}

	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, isOpen := <-input:
			if !isOpen {
				return nil
			}

			if strings.Contains(value, decoratorSkipMessage) {
				return ErrCannotDecorate
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
		for _, ch := range outputs {
			close(ch)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	current := 0
	total := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, isOpen := <-input:
			if !isOpen {
				return nil
			}

			select {
			case outputs[current] <- value:
				current = (current + 1) % total
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

	processInput := func(ch chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case value, isOpen := <-ch:
				if !isOpen {
					return
				}

				if strings.Contains(value, multiplexerSkip) {
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
		go processInput(ch)
	}

	wg.Wait()
	return nil
}