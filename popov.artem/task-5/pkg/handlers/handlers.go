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
		case item, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(item, "no decorator") {
				return fmt.Errorf("can’t be decorated")
			}
			if !strings.HasPrefix(item, "decorated:") {
				item = "decorated:" + item
			}
			select {
			case output <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				return nil
			}
			outputIndex := counter % len(outputs)
			select {
			case outputs[outputIndex] <- item:
			case <-ctx.Done():
				return nil
			}
			counter++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	multiplexedChan := make(chan string)

	var multiplexerWg sync.WaitGroup
	multiplexerWg.Add(len(inputs))

	for _, inCh := range inputs {
		go func(ch chan string) {
			defer multiplexerWg.Done()
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
		multiplexerWg.Wait()
		close(multiplexedChan)
	}()

	for item := range multiplexedChan {
		select {
		case output <- item:
		case <-ctx.Done():
			return nil
		}
	}

	return nil
}
