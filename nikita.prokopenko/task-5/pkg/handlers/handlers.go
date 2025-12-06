package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	noDecorator = "no decorator"
	decorated   = "decorated: "
	noMux       = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	if output == nil {
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

	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, noDecorator) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(val, decorated) {
				val = decorated + val
			}

			select {
			case output <- val:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		closed := make(map[chan string]bool)
		for _, outCh := range outputs {
			if outCh == nil {
				continue
			}
			if !closed[outCh] {
				close(outCh)
				closed[outCh] = true
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
		case val, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[idx] <- val:
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

	if len(inputs) == 0 {
		close(output)
		return nil
	}

	var wg sync.WaitGroup
	var once sync.Once
	defer func() {
		wg.Wait()
		once.Do(func() { close(output) })
	}()

	for _, ch := range inputs {
		wg.Add(1)
		go func(in chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(val, noMux) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- val:
					}
				}
			}
		}(ch)
	}

	return nil
}