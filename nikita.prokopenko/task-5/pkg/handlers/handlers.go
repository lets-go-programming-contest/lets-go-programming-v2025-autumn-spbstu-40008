// pkg/handlers/handlers.go
package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrDecoratingImpossible = errors.New("can't be decorated")

const (
	skipDecoration     = "no decorator"
	decorationPrefix   = "decorated: "
	skipMultiplexing   = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(item, skipDecoration) {
				return ErrDecoratingImpossible
			}

			if !strings.HasPrefix(item, decorationPrefix) {
				item = decorationPrefix + item
			}

			select {
			case output <- item:
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
		case item, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[current] <- item:
				current = (current + 1) % total
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

	processStream := func(stream chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-stream:
				if !ok {
					return
				}

				if strings.Contains(item, skipMultiplexing) {
					continue
				}

				select {
				case output <- item:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, stream := range inputs {
		wg.Add(1)
		go processStream(stream)
	}

	wg.Wait()
	return nil
}