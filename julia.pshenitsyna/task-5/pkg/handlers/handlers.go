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
		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(val, "decorated:") {
				val = "decorated:" + val
			}
			select {
			case <-ctx.Done():
				return nil
			case output <- val:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0{
		return nil
	}
	
	cnt := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
			if !ok {
				return nil
			}

			outputChan := outputs[cnt]
			if outputChan == nil {
				cnt++
				if cnt >= len(outputs) {
					cnt = 0
				}
				continue
			}

			select {
			case <-ctx.Done():

				return nil

			case outputChan <- val:
			}

			cnt++
			if cnt >= len(outputs) {
				cnt = 0
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	for _, inputChan := range inputs {
		waitGroup.Add(1)

		go func(channel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-channel:
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
		}(inputChan)
	}
	waitGroup.Wait()

	return nil
}
