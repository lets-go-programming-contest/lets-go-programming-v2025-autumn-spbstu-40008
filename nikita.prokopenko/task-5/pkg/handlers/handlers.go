package handlers

import (
	"context"
	"strings"
)

type DecoratorError string

func (e DecoratorError) Error() string {
	return string(e)
}

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "
	if output == nil {
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

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return DecoratorError("can't be decorated")
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
		return nil
	}

	doneChan := make(chan struct{})
	defer close(doneChan)

	type workerResult struct{}
	workerCount := len(inputs)

	for _, inChan := range inputs {
		ch := inChan

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						workerCount--
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
		}()
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if workerCount <= 0 {
				return nil
			}
		}
	}
}
