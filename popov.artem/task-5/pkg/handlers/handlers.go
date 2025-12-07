package handlers

import (
	"context"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(item, "no decorator") {
				return context.Canceled
			}
			if !strings.HasPrefix(item, "decorated:") {
				item = "decorated:" + item
			}
			select {
			case output <- item:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	counter := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-input:
			if !ok {
				return nil
			}
			outputIndex := counter % len(outputs)
			select {
			case outputs[outputIndex] <- item:
			case <-ctx.Done():
				return ctx.Err()
			}
			counter++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	multiplexedChan := make(chan string)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, inCh := range inputs {
		go func(ch chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-ch:
					if !ok {
						return
					}
					if !strings.Contains(item, "no multiplexer") {
						select {
						case multiplexedChan <- item:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}(inCh)
	}

	go func() {
		wg.Wait()
		close(multiplexedChan)
	}()

	for item := range multiplexedChan {
		select {
		case output <- item:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}