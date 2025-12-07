package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can’t be decorated")

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
				return ErrCannotBeDecorated
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
	var multiplexerWaitGroup sync.WaitGroup
	multiplexerWaitGroup.Add(len(inputs))

	for _, inputChannel := range inputs {
		go func(currentInput chan string) {
			defer multiplexerWaitGroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-currentInput:
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
		multiplexerWaitGroup.Wait()
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
