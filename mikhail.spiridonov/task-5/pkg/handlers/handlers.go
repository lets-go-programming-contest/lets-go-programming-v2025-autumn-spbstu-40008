package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

	const errorSubstring = "no decorator"

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, errorSubstring) {
				return fmt.Errorf("%w: %s", ErrCantBeDecorated, data)
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
	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := counter % len(outputs)
			counter++

			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	const skipSubstring = "no multiplexer"

	var workerGroup sync.WaitGroup

	workerErrors := make(chan error, len(inputs))

	processChannel := func(inputChan chan string) {
		defer workerGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-inputChan:
				if !ok {
					return
				}

				if strings.Contains(data, skipSubstring) {
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
		workerGroup.Add(1)

		go processChannel(inputChan)
	}

	done := make(chan struct{})
	go func() {
		workerGroup.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		select {
		case <-done:
			return nil
		case err := <-workerErrors:
			return err
		}
	case err := <-workerErrors:
		return err
	}
}
