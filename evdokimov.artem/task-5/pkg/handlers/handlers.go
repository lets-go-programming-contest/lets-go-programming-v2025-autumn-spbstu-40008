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

func PrefixDecoratorFunc(ctx context.Context, input <-chan string, output chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
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
			case output <- val:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input <-chan string, outputs []chan<- string) error {
	if len(outputs) == 0 {
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
			case <-ctx.Done():
				return nil
			}
			idx = (idx + 1) % len(outputs)
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []<-chan string, output chan<- string) error {
	var waitGroup sync.WaitGroup
	for _, inCh := range inputs {
		channel := inCh
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
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
					case output <- val:
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
