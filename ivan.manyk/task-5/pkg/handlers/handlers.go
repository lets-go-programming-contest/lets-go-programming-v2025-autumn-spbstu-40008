package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrDecorationNotAllowed = errors.New("can't be decorated")

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
				return ErrDecorationNotAllowed
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := counter % len(outputs)
			counter++

			select {
			case <-ctx.Done():
				return nil
			case outputs[idx] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	done := make(chan string)

	for _, inputCh := range inputs {
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
					case <-ctx.Done():
						return
					case done <- data:
					}
				}
			}
		}(inputCh)
	}

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-done:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}
