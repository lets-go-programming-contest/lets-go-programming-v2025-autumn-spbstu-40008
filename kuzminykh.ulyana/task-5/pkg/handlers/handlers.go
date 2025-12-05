package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	const prefix = "decorated: "
	for {
		select {
		case data, ok := <-input:
			if !ok {
				close(output)
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- data:
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		close(output)
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	done := make(chan struct{})
	var doneOnce sync.Once

	go func() {
		wg.Wait()
		doneOnce.Do(func() {
			close(output)
			close(done)
		})
	}()

	for _, in := range inputs {
		go func(in <-chan string) {
			defer wg.Done()
			for {
				select {
				case data, ok := <-in:
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
				case <-ctx.Done():
					return
				}
			}
		}(in)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var counter int64 = 0

	for {
		select {
		case data, ok := <-input:
			if !ok {
				for _, out := range outputs {
					close(out)
				}
				return nil
			}

			idx := atomic.AddInt64(&counter, 1) - 1
			ch := outputs[int(idx)%len(outputs)]

			select {
			case <-ctx.Done():
				return ctx.Err()
			case ch <- data:
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
