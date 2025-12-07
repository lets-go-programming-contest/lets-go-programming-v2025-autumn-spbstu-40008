package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	const prefix = "decorated: "

	for {
		select {
		case data, exist := <-input:
			if !exist {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}
			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	done := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	for _, inputChan := range inputs {
		go func(in <-chan string) {
			defer waitGroup.Done()

			for {
				select {
				case data, exist := <-in:
					if !exist {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				case <-ctx.Done():
					return
				}
			}
		}(inputChan)
	}

	select {
	case <-ctx.Done():
		<-done
		return nil
	case <-done:
		return nil
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, exist := <-input:
				if !exist {
					return nil
				}
			}
		}
	}

	var counter int64 = 0

	for {
		select {
		case data, exist := <-input:
			if !exist {
				return nil
			}

			idx := atomic.AddInt64(&counter, 1) - 1
			outChan := outputs[int(idx)%len(outputs)]

			select {
			case <-ctx.Done():
				return nil
			case outChan <- data:
			}
		case <-ctx.Done():
			return nil
		}
	}
}
