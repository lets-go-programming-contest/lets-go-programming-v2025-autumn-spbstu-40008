package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	prefix := "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil

		case item, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(item, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(item, prefix) {
				item = prefix + item
			}

			select {
			case output <- item:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, output []chan string) error {
	counter := 0
	outputCount := len(output)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case item, ok := <-input:
			if !ok {
				return nil
			}

			targetChanIndex := counter % outputCount
			counter++

			targetChan := output[targetChanIndex]

			select {
			case targetChan <- item:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup

	merge := func(ch chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return

			case item, ok := <-ch:
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
	}

	for _, ch := range inputs {
		wg.Add(1)
		go merge(ch)
	}

	wg.Wait()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
