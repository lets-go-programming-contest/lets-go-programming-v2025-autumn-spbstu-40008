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

			sent := false
			for attempts := 0; attempts < len(outputs) && !sent; attempts++ {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case outputs[i] <- data:
					sent = true
					i = (i + 1) % len(outputs)
				default:
					i = (i + 1) % len(outputs)
				}
			}
		}
	}
}
