package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(
	ctx context.Context,
	inputChan, outputChan chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inputChan:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case outputChan <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	inputChan chan string,
	outputChans []chan string,
) error {
	if len(outputChans) == 0 {
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, ok := <-inputChan:
				if !ok {
					return nil
				}
			}
		}
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inputChan:
			if !ok {
				return nil
			}

			select {
			case outputChans[index] <- data:
			case <-ctx.Done():
				return nil
			}

			index = (index + 1) % len(outputChans)
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputChans []chan string,
	outputChan chan string,
) error {
	if len(inputChans) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	for _, inputChannel := range inputChans {
		waitGroup.Add(1)

		go func(ch chan string) {
			defer waitGroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return
					}
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case outputChan <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inputChannel)
	}

	done := make(chan struct{})
	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		waitGroup.Wait()
		return nil
	}
}