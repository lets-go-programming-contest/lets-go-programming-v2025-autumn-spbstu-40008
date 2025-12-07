package handlers

import (
	"context"
	"strings"
	"sync"
	"time"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup
	done := make(chan struct{})

	for _, in := range inputs {
		wg.Add(1)
		go func(in chan string) {
			defer wg.Done()

			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					case <-done:
						return
					case <-time.After(50 * time.Millisecond):
						continue
					}
				case <-time.After(100 * time.Millisecond):
					continue
				}
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
