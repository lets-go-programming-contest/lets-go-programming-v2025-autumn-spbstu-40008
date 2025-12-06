package handlers

import (
	"context"
	"strings"
	"sync"
)

type DecoratorError string

func (e DecoratorError) Error() string {
	return string(e)
}

func drainInput(ctx context.Context, inputChan chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case _, ok := <-inputChan:
			if !ok {
				return nil
			}
		}
	}
}

func doPrefixDecorator(ctx context.Context, inputChan, outputChan chan string) error {
	const prefix = "decorated: "

	defer close(outputChan)

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-inputChan:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return DecoratorError("can't be decorated")
			}

			if !strings.HasPrefix(value, prefix) {
				value = prefix + value
			}

			select {
			case outputChan <- value:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func PrefixDecoratorFunc(ctx context.Context, inputChan, outputChan chan string) error {
	if outputChan == nil {
		return drainInput(ctx, inputChan)
	}
	return doPrefixDecorator(ctx, inputChan, outputChan)
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return drainInput(ctx, input)
	}

	index := 0

	defer func() {
		for _, out := range outputs {
			select {
			case _, ok := <-out:
				if ok {
				}
			default:
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- value:
				index = (index + 1) % len(outputs)
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if output == nil {
		return nil
	}

	if len(inputs) == 0 {
		defer close(output)
		return nil
	}

	var waitGroup sync.WaitGroup

	for _, inputCh := range inputs {
		source := inputCh
		waitGroup.Add(1)

		go func(src chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-src:
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
		}(source)
	}

	waitGroup.Wait()
	defer close(output)

	return nil
}
