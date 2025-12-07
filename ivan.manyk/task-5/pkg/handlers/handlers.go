package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return fmt.Errorf("can't be decorated")
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := counter % len(outputs)
			counter++

			select {
			case <-ctx.Done():
				return nil
			case outputs[idx] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	aggregate := make(chan string)

	for _, inputCh := range inputs {
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
					select {
					case <-ctx.Done():
						return
					case aggregate <- data:
					}
				}
			}
		}(inputCh)
	}

	go func() {
		wg.Wait()
		close(aggregate)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-aggregate:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}
