package handlers

import (
	"context"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
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

	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(value, "no decorator") {
				return DecoratorError("can't be decorated")
			}
			if !strings.HasPrefix(value, "decorated: ") {
				value = "decorated: " + value
			}
			select {
			case output <- value:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		closed := make(map[chan string]bool)
		for _, out := range outputs {
			if out == nil {
				continue
			}
			if !closed[out] {
				close(out)
				closed[out] = true
			}
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case value, ok := <-input:
			if !ok {
				return nil
			}
			select {
			case outputs[idx] <- value:
				idx = (idx + 1) % len(outputs)
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

	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	done := make(chan bool, len(inputs))

	for _, in := range inputs {
		go func(ch chan string) {
			for {
				select {
				case <-ctx.Done():
					done <- true
					return
				case val, ok := <-ch:
					if !ok {
						done <- true
						return
					}
					if strings.Contains(val, "no multiplexer") {
						continue
					}
					select {
					case <-ctx.Done():
						done <- true
						return
					case output <- val:
					}
				}
			}
		}(in)
	}

	for i := 0; i < len(inputs); i++ {
		<-done
	}

	return nil
}

type DecoratorError string

func (e DecoratorError) Error() string {
	return string(e)
}
