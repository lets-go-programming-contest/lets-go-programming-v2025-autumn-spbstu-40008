package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	prefix := "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- data:
			case <-ctx.Done():
				return nil
			}

			index = (index + 1) % len(outputs)
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	processInput := func(inputChan chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case data, ok := <-inputChan:
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
				}
			}
		}
	}

	for _, inputChan := range inputs {
		waitGroup.Add(1)

		go processInput(inputChan)
	}

	waitGroup.Wait()

	return nil
}
