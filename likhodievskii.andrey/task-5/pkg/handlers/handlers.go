package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoDecorator  = errors.New("can't be decorated")
	ErrEmptyOutputs = errors.New("empty outputs")
)

const (
	noDecorator     = "no decorator"
	noMultiplexer   = "no multiplexer"
	decoratedPrefix = "decorated: "
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case str, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(str, noDecorator) {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(str, decoratedPrefix) {
				str = decoratedPrefix + str
			}

			select {
			case output <- str:

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

	index := 0

	for {
		select {
		case str, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- str:

			case <-ctx.Done():
				return nil
			}

			index = (index + 1) % len(outputs)

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
	waitgr := sync.WaitGroup{}

	waitgr.Add(len(inputs))

	for _, channel := range inputs {
		go func() {
			defer waitgr.Done()

			for {
				select {
				case str, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(str, noMultiplexer) {
						continue
					}

					select {
					case output <- str:

					case <-ctx.Done():
						return
					}

				case <-ctx.Done():
					return
				}
			}
		}()
	}

	waitgr.Wait()

	return nil
}
