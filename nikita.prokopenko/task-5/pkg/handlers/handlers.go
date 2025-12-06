// pkg/handlers/handlers.go
package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, in, out chan string) error {
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-in:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
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
	defer func() {
		for _, ch := range outs {
			close(ch)
		}
	}()

	if len(outs) == 0 {
		return nil
	}

	idx := 0
	total := len(outs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-in:
			if !ok {
				return nil
			}

			select {
			case outs[idx] <- val:
				idx = (idx + 1) % total
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, ins []chan string, out chan string) error {
	defer close(out)

	if len(ins) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	process := func(ch chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch:
				if !ok {
					return
				}

				if strings.Contains(val, "no multiplexer") {
					continue
				}

				select {
				case out <- val:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, ch := range ins {
		wg.Add(1)
		go process(ch)
	}

	wg.Wait()
	return nil
}