package handlers

import (
	"context"
	"strings"
)

func drainInput(ctx context.Context, input chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
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
			return ctx.Err()
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
				return ctx.Err()
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
			return ctx.Err()
		case value, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[idx] <- value:
				idx = (idx + 1) % len(outputs)
			case <-ctx.Done():
				return ctx.Err()
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

	doneCh := make(chan struct{}, len(inputs))

	for _, inputCh := range inputs {
		go func(inputCh chan string) {
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
		}(inputCh)
	}

	for i := 0; i < len(inputs); i++ {
		select {
		case <-doneCh:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

type DecoratorError string

func (e DecoratorError) Error() string {
	return string(e)
}