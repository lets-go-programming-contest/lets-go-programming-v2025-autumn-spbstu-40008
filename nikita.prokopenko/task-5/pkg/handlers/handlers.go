package handlers

import (
	"context"
	"errors"
	"strings"
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
			case <-ctx.Done():
				return nil
			case output <- val:
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
	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[idx] <- val:
				idx = (idx + 1) % len(outputs)
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

	doneCh := make(chan struct{})
	defer close(doneCh)

	for _, ch := range inputs {
		go func(in chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-doneCh:
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
					case <-doneCh:
						return
					case output <- val:
					}
				}
			}
		}(ch)
	}

	<-ctx.Done()
	return nil
}