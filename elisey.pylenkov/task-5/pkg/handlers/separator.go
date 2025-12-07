package handlers

import (
	"context"
)

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}
			for sent := false; !sent; {
				ch := outputs[i%len(outputs)]
				select {
				case ch <- data:
					sent = true
				case <-ctx.Done():
					return ctx.Err()
				default:
				}
				i++
			}
		}
	}
}
