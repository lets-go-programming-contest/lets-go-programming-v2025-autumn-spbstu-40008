package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

func PrefixDecorator(ctx context.Context, input chan string, output chan string) error {
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
			case <-ctx.Done():
				return nil
			case output <- item:
			}
		}
	}
}

func Multiplexer(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup
	for _, ch := range inputs {
		wg.Add(1)
		go func(in chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case item, ok := <-in:
					if !ok {
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
		}(ch)
	}
	wg.Wait()

	select {
	case <-ctx.Done():
		return nil
	default:
		close(output)
	}
	return nil
}

func Separator(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0
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