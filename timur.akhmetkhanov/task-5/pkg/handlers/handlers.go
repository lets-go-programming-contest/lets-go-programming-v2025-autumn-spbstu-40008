package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

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
				return errors.New("can't be decorated")
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
	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-input:
			if !ok {
				for _, ch := range outputs {
					close(ch)
				}
				return nil
			}

			targetCh := outputs[i%len(outputs)]
			i++

			select {
			case targetCh <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	for _, ch := range inputs {
		wg.Add(1)
		go func(c chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-c:
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
		}(ch)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	<-ctx.Done()
	return nil
}
