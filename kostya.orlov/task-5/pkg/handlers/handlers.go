package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecoratorRequiresOutput   = errors.New("decorator requires at least one output channel")
	ErrSeparatorRequiresOutput   = errors.New("separator requires at least one output channel")
	ErrMultiplexerRequiresOutput = errors.New("multiplexer requires at least one output channel")
	ErrCantBeDecorated           = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrDecoratorRequiresOutput
	}
	output := outputs[0]

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
	if len(outputs) == 0 {
		return ErrSeparatorRequiresOutput
	}

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

func MultiplexerFunc(ctx context.Context, inputs []chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrMultiplexerRequiresOutput
	}
	output := outputs[0]

	var waitGroup sync.WaitGroup
	fanIn := make(chan string)

	for _, inputChannel := range inputs {
		waitGroup.Add(1)
		go func(channel chan string) {
			defer waitGroup.Done()
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
		waitGroup.Wait()
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