package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrDecorateFail = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, in chan string, out chan string) error {
	const badPhrase = "no decorator"
	const pre = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-in:
			if !ok {
				return nil
			}
			if strings.Contains(val, badPhrase) {
				return ErrDecorateFail
			}
			if !strings.HasPrefix(val, pre) {
				val = pre + val
			}
			select {
			case out <- val:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, in chan string, outs []chan string) error {
	n := len(outs)
	if n == 0 {
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, ok := <-in:
				if !ok {
					return nil
				}
			}
		}
	}

	pos := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-in:
			if !ok {
				return nil
			}
			dest := outs[pos]
			select {
			case dest <- item:
				pos = (pos + 1) % n
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, sources []chan string, sink chan string) error {
	if len(sources) == 0 {
		return nil
	}

	var group sync.WaitGroup
	group.Add(len(sources))

	for _, src := range sources {
		go func(ch chan string) {
			defer group.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-ch:
					if !ok {
						return
					}
					if strings.Contains(msg, "no multiplexer") {
						continue
					}
					select {
					case sink <- msg:
					case <-ctx.Done():
						return
					}
				}
			}
		}(src)
	}

	group.Wait()
	return nil
}
