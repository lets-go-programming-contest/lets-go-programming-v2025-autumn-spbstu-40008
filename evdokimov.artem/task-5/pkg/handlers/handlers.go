package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

const (
	NoDecoratorKey   = "no decorator"
	DecoratedPrefix  = "decorated: "
	NoMultiplexerKey = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, inputChannel chan string, outputChannel chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case val, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(val, NoDecoratorKey) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(val, DecoratedPrefix) {
				val = DecoratedPrefix + val
			}

			select {
			case outputChannel <- val:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, inputChannel chan string, outputChannels []chan string) error {
	if len(outputChannels) == 0 {
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, ok := <-inputChannel:
				if !ok {
					return nil
				}
			}
		}
	}

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case val, ok := <-inputChannel:
			if !ok {
				return nil
			}

			select {
			case outputChannels[idx] <- val:
			case <-ctx.Done():
				return nil
			}

			idx = (idx + 1) % len(outputChannels)
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	var wg sync.WaitGroup
	wg.Add(len(inputChannels))

	for _, inCh := range inputChannels {
		channel := inCh

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case val, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(val, NoMultiplexerKey) {
						continue
					}

					select {
					case outputChannel <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	wg.Wait()

	return nil
}
