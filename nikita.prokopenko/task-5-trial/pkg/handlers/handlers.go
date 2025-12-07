// pkg/handlers/handlers.go
package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrDecorateFail = errors.New("can't be decorated")

const (
	stopWord = "no decorator"
	prefix   = "decorated: "
)

// drainInput просто читает вход и ничего не делает — используется когда нет куда писать.
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

// doPrefixDecorator содержит основную логику декоратора.
func doPrefixDecorator(ctx context.Context, inputChan, outputChan chan string) error {
	defer close(outputChan)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-inputChan:
			if !ok {
				return nil
			}

			if strings.Contains(msg, stopWord) {
				return ErrDecorateFail
			}

			if !strings.HasPrefix(msg, prefix) {
				msg = prefix + msg
			}

			select {
			case outputChan <- msg:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

// PrefixDecoratorFunc — тонкая обёртка, чтобы привести поведение к линтеру (простая функция).
func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	if output == nil {
		return drainInput(ctx, input)
	}
	return doPrefixDecorator(ctx, input, output)
}

// SeparatorFunc распределяет сообщения по выходным каналам по очереди.
// Если outputs пустой — просто сливаeт вход.
func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return drainInput(ctx, input)
	}

	index := 0
	count := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			target := outputs[index]

			select {
			case target <- msg:
				index = (index + 1) % count
			case <-ctx.Done():
				return nil
			}
		}
	}
}

// MultiplexerFunc объединяет несколько входных каналов в один выходной.
// Для линтера вынес worker как именованную переменную и вызывает её через go worker(pipe).
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	worker := func(src chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-src:
				if !ok {
					return
				}

				if strings.Contains(msg, "no multiplexer") {
					continue
				}

				select {
				case output <- msg:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, inputPipe := range inputs {
		waitGroup.Add(1)
		go worker(inputPipe)
	}

	waitGroup.Wait()

	return nil
}
