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
	in chan string,
	out chan string,
) error {
	for {
		select {
		case data, ok := <-in:
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
			case out <- data:
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
	in chan string,
	outs []chan string,
) error {
	if len(outs) == 0 {
		return ErrEmptyOutputs
	}
	idx := 0
	for {
		select {
		case data, ok := <-in:
			if !ok {
				return nil
			}
			select {
			case outs[idx] <- data:
			case <-ctx.Done():
				return nil
			}
			idx = (idx + 1) % len(outs)
		case <-ctx.Done():
			return nil
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	ins []chan string,
	out chan string,
) error {
	var wg sync.WaitGroup
	wg.Add(len(ins))

	for _, ch := range ins {
		ch := ch
		go func() {
			defer wg.Done()
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
					case out <- data:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
