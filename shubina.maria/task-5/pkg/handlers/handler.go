package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

// PrefixDecoratorFunc — модификатор, добавляющий префикс.
func PrefixDecoratorFunc(
	ctx context.Context,
	input, output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
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
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

// SeparatorFunc — распределяет входящие данные по выходным каналам по кругу.
func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		// Просто сливаем вход, если нет выходов
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, ok := <-input:
				if !ok {
					return nil
				}
			}
		}
	}

	idx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return nil
			}
			idx = (idx + 1) % len(outputs)
		}
	}
}

// MultiplexerFunc — объединяет несколько входов в один выход, фильтруя "no multiplexer".
func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	for _, ch := range inputs {
		wg.Add(1)
		go func(in chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(ch)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		wg.Wait()
		return nil
	}
}