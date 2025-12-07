package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, isOpen := <-input:
			if !isOpen {
				close(output)
				return nil
			}

			if strings.Contains(item, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(item, "decorated: ") {
				item = "decorated: " + item
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- item:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	for _, inputChan := range inputs {
		waitGroup.Add(1)

		go func(inCh chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case item, isOpen := <-inCh:
					if !isOpen {
						return
					}

					if strings.Contains(item, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- item:
					}
				}
			}
		}(inputChan)
	}

	waitGroup.Wait()

	select {
	case <-ctx.Done():
		return nil
	default:
		close(output)
	}

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, isOpen := <-input:
			if !isOpen {
				for _, outCh := range outputs {
					close(outCh)
				}
				return nil
			}

			index := counter % len(outputs)
			targetChan := outputs[index]
			counter++

			select {
			case <-ctx.Done():
				return nil
			case targetChan <- item:
			}
		}
	}
}