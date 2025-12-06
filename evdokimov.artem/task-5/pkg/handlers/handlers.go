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
	NoMultiplexerKey = "no multiplexer"
	DecoratedPrefix  = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, inputChannel chan string, outputChannel chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-inputChannel:
			if !ok {
				return nil
			}
			if strings.Contains(value, NoDecoratorKey) {
				return ErrCantBeDecorated
			}
			if !strings.HasPrefix(value, DecoratedPrefix) {
				value = DecoratedPrefix + value
			}
			select {
			case outputChannel <- value:
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
			case _, ok := <-inputChannel:
				if !ok {
					return nil
				}
			case <-ctx.Done():
				return nil
			}
		}
	}

	idx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-inputChannel:
			if !ok {
				return nil
			}
			select {
			case outputChannels[idx] <- value:
			case <-ctx.Done():
				return nil
			}
			idx = (idx + 1) % len(outputChannels)
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(inputChannels))

	for _, ch := range inputChannels {
		channel := ch
		go func() {
			defer waitGroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-channel:
					if !ok {
						return
					}
					if strings.Contains(value, NoMultiplexerKey) {
						continue
					}
					select {
					case outputChannel <- value:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	waitGroup.Wait()
	return nil
}
