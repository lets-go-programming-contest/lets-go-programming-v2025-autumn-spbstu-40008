package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrInvalidContent = errors.New("can't be decorated")

const (
	skipDecorator   = "no decorator"
	skipMultiplexer = "no multiplexer"
	prefix          = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, skipDecorator) {
				return ErrInvalidContent
			}

			res := val
			if !strings.HasPrefix(val, prefix) {
				res = prefix + val
			}

			select {
			case output <- res:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var idx int
	count := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
			if !ok {
				return nil
			}

			target := outputs[idx%count]
			idx++

			select {
			case target <- val:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	handle := func(ch chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch:
				if !ok {
					return
				}

				if strings.Contains(val, skipMultiplexer) {
					continue
				}

				select {
				case output <- val:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, ch := range inputs {
		wg.Add(1)

		go handle(ch)
	}

	wg.Wait()
	return nil
}
