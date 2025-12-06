package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

func drainInput(ctx context.Context, input chan string) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during drain: %w", ctx.Err())
		case _, ok := <-input:
			if !ok {
				return nil
			}
		}
	}
}

func runPrefixDecorator(ctx context.Context, input, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during decoration: %w", ctx.Err())
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
				return fmt.Errorf("context cancelled during sending decorated value: %w", ctx.Err())
			}
		}
	}
}

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	if output == nil {
		return drainInput(ctx, input)
	}

	return runPrefixDecorator(ctx, input, output)
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		closed := make(map[chan string]bool)
		for _, outCh := range outputs {
			if outCh == nil {
				continue
			}

			if !closed[outCh] {
				close(outCh)
				closed[outCh] = true
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
			return fmt.Errorf("context cancelled during separation: %w", ctx.Err())
		case value, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[idx] <- value:
				idx = (idx + 1) % len(outputs)
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during sending separated value: %w", ctx.Err())
			}
		}
	}
}

func processInputChannel(ctx context.Context, inputCh chan string, output chan string, doneCh chan struct{}) {
	defer func() {
		doneCh <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-inputCh:
			if !ok {
				return
			}

			if strings.Contains(val, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case output <- val:
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

	var wg sync.WaitGroup
	doneCh := make(chan struct{}, len(inputs))

	for _, inputCh := range inputs {
		wg.Add(1)
		go func(ch chan string) {
			defer wg.Done()
			processInputChannel(ctx, ch, output, doneCh)
		}(inputCh)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for range inputs {
		select {
		case <-doneCh:
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during multiplexing: %w", ctx.Err())
		}
	}

	return nil
}

type DecoratorError string

func (e DecoratorError) Error() string {
	return string(e)
}

func (e DecoratorError) Unwrap() error {
	return errors.New(string(e))
}