package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
	ErrEmptyOutputs    = errors.New("empty outputs")
)

const (
	noDecorator     = "no decorator"
	noMultiplexer   = "no multiplexer"
	decoratedPrefix = "decorated: "
)

func PrefixDecoratorFunc(
	ctx context.Context,
	inputChannel chan string,
	outputChannel chan string,
) error {
	for {
		select {
		case data, ok := <-inputChannel:
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
			case outputChannel <- data:
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
	inputChannel chan string,
	outputChannels []chan string,
) error {
	if len(outputChannels) == 0 {
		return ErrEmptyOutputs
	}

	index := 0

	for {
		select {
		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			select {
			case outputChannels[index] <- data:
			case <-ctx.Done():
				return nil
			}

			index = (index + 1) % len(outputChannels)
		case <-ctx.Done():
			return nil
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputChannels []chan string,
	outputChannel chan string,
) error {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputChannels))

	for _, inputChannel := range inputChannels {
		go func(ch chan string) {
			defer waitGroup.Done()

			for {
				select {
				case data, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(data, noMultiplexer) {
						continue
					}

					select {
					case outputChannel <- data:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}(inputChannel)
	}

	waitGroup.Wait()

	return nil
}
