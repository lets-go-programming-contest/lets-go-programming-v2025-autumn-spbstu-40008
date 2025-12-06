package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
	ErrEmptyOutputs    = errors.New("empty outputs list")
)

const (
	noDecorator     = "no decorator"
	noMultiplexer   = "no multiplexer"
	decoratedPrefix = "decorated: "
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input, output chan string,
) error {
	for {
		select {
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecorator) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, decoratedPrefix) {
				data = decoratedPrefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	idx := 0

	for {
		select {
		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return nil
			}

			idx = (idx + 1) % len(outputs)

		case <-ctx.Done():
			return nil
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, input := range inputs {

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case data, ok := <-input:
					if !ok {
						return
					}

					if strings.Contains(data, noMultiplexer) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}

				case <-ctx.Done():
					return
				}
			}
		}()
	}

	waitGroup.Wait()

	return nil
}
