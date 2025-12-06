package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrDecorateFail = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, in <-chan string, out chan<- string) error {
	const prefix = "decorated: "
	const stopWord = "no decorator"

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-in:
			if !ok {
				return nil 
			}

			if strings.Contains(msg, stopWord) {
				return ErrDecorateFail
			}

			if !strings.HasPrefix(msg, prefix) {
				msg = prefix + msg
			}

			select {
			case out <- msg:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, in <-chan string, outs []chan string) error {
	if len(outs) == 0 {
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

	idx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-in:
			if !ok {
				return nil
			}

			target := outs[idx]
			
			select {
			case target <- msg:
				idx = (idx + 1) % len(outs)
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, ins []chan string, out chan<- string) error {
	var wg sync.WaitGroup

	mergeRoutine := func(ch <-chan string) {
		defer wg.Done()
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
				case out <- msg:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, ch := range ins {
		wg.Add(1)
		go mergeRoutine(ch)
	}

	wg.Wait()
	
	return nil
}