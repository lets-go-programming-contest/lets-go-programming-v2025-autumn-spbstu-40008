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
			return nil

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
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {

	if len(outputs) == 0 {
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, ok := <-input:
				if !ok {
					return nil
				}
			}
		}
	}

	index := 0
	numOutputs := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return nil

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
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	fanIn := make(chan string, 100)
	var wg sync.WaitGroup

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
					case fanIn <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inputCh)
	}

	go func() {
		wg.Wait()
		close(fanIn)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
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
				return nil
			}
		}
	}
}