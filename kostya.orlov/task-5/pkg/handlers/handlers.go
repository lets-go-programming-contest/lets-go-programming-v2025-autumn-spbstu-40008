package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, "decorated:") {
				data = "decorated:" + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	index := 0
	numOutputs := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case data, ok := <-input:
			if !ok {
				return nil
			}

			outputIndex := index % numOutputs
			targetOutput := outputs[outputIndex]

			select {
			case targetOutput <- data:
				index++
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup
	fanIn := make(chan string)

	for _, inputChannel := range inputs {
		wg.Add(1)
		go func(channel chan string) {
			defer wg.Done()
			
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-channel:
					if !ok {
						return
					}
					select {
					case fanIn <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inputChannel)
	}

	go func() {
		wg.Wait()
		close(fanIn)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case data, ok := <-fanIn:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no multiplexer") {
				continue
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}