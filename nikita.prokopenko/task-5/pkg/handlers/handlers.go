// pkg/handlers/handlers.go
package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrDecorationImpossible = errors.New("can't be decorated")

func ApplyPrefix(ctx context.Context, source, dest chan string) error {
	prefix := "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-source:
			if !ok {
				return nil
			}

			if strings.Contains(item, "no decorator") {
				return ErrDecorationImpossible
			}

			if !strings.HasPrefix(item, prefix) {
				item = prefix + item
			}

			select {
			case dest <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func Distribute(ctx context.Context, source chan string, targets []chan string) error {
	if len(targets) == 0 {
		// No targets to distribute to, just drain the source
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, ok := <-source:
				if !ok {
					return nil
				}
			}
		}
	}

	current := 0
	total := len(targets)

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-source:
			if !ok {
				return nil
			}

			select {
			case targets[current] <- item:
				current = (current + 1) % total
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func Merge(ctx context.Context, sources []chan string, result chan string) error {
	if result == nil || len(sources) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	processSource := func(ch chan string) {
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
				case result <- item:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, ch := range sources {
		wg.Add(1)
		go processSource(ch)
	}

	wg.Wait()
	return nil
}