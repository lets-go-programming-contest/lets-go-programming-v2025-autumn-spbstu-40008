package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrDecorateFail = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

	const stopWord = "no decorator" 

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-input:
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
			case output <- msg:
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

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			target := outputs[idx]

			select {
			case target <- msg:
				idx = (idx + 1) % len(outputs)
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	if len(inputs) == 0 {
		return nil
	}

	mergeRoutine := func(inputChannel chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-inputChannel:
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

	for _, channel := range inputs {
		waitGroup.Add(1)

		go mergeRoutine(channel)
	}

	waitGroup.Wait()

	return nil
}
