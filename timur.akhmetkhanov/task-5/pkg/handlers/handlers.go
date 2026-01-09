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
		case item, ok := <-input:
			if !ok {
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
			case output <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				for _, channel := range outputs {
					close(channel)
				}

				return nil
			}

			targetCh := outputs[index%len(outputs)]
			index++

			select {
			case targetCh <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	for _, channel := range inputs {
		waitGroup.Add(1)

		go func(workerChan chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-workerChan:
					if !ok {
						return
					}

					if strings.Contains(item, "no multiplexer") {
						continue
					}

					select {
					case output <- item:
					case <-ctx.Done():
						return
					}
				}
			}
		}(channel)
	}

	go func() {
		waitGroup.Wait()
		close(output)
	}()

	<-ctx.Done()

	return nil
}
