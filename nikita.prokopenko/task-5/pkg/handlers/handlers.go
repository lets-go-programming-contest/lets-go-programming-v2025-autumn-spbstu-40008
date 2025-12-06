package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	noDecoratorMessage   = "no decorator"
	decoratorPrefix      = "decorated: "
	noMultiplexerMessage = "no multiplexer"
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

			if strings.Contains(value, noDecoratorMessage) {
				return ErrCannotBeDecorated
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
		for _, outCh := range outputs {
			if outCh != nil {
				close(outCh)
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
	if output == nil {
		return nil
	}

	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	defer close(done)

	for _, inputCh := range inputs {
		wg.Add(1)
		go func(ch chan string) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case val, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(val, noMultiplexerMessage) {
						continue
					}

					select {
					case <-done:
						return
					case output <- val:
					}
				}
			}
		}(inputCh)
	}

	go func() {
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return nil
	}
}

type DecoratorError string

func (e DecoratorError) Error() string {
	return string(e)
}