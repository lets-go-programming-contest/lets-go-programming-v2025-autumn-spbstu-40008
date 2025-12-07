package handlers

import (
	"context"
	"strings"
	"sync"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)

		go func(ch chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return

				case data, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}(in)
	}

	wg.Wait()
	return nil
}
