package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotDecorate = fmt.Errorf("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled in PrefixDecoratorFunc: %w", ctx.Err())
		case item, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(item, "no decorator") {
				return fmt.Errorf("%w: %s", ErrCannotDecorate, item)
			}
			if !strings.HasPrefix(item, "decorated:") {
				item = "decorated:" + item
			}
			select {
			case output <- item:
			case <-ctx.Done():
				return fmt.Errorf("context canceled while sending in PrefixDecoratorFunc: %w", ctx.Err())
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, outputChannel := range outputs {
			close(outputChannel)
		}
	}()

	var counter int
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled in SeparatorFunc: %w", ctx.Err())
		case item, ok := <-input:
			if !ok {
				return nil
			}
			outputIndex := counter % len(outputs)
			select {
			case outputs[outputIndex] <- item:
			case <-ctx.Done():
				return fmt.Errorf("context canceled while sending in SeparatorFunc: %w", ctx.Err())
			}
			counter++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	multiplexedChan := make(chan string)

	var multiplexerWg sync.WaitGroup
	multiplexerWg.Add(len(inputs))

	for _, inputChannel := range inputs {
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
		}(inputChannel)
	}

	go func() {
		multiplexerWg.Wait()
		close(multiplexedChan)
	}()

	for item := range multiplexedChan {
		select {
		case output <- item:
		case <-ctx.Done():
			return fmt.Errorf("context canceled while sending in MultiplexerFunc: %w", ctx.Err())
		}
	}

	return nil
}
