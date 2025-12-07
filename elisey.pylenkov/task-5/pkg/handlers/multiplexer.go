package handlers

import (
	"context"
	"sync"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(in chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(in)
	}

	wg.Wait()
	return nil
}
