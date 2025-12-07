package handlers

import (
	"context"
	"strings"
)

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for _, in := range inputs {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case data, ok := <-in:
					if !ok {
						continue
					}
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					case output <- data:
					}
				default:
				}
			}
		}
	}
}
