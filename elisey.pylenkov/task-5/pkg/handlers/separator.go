package handlers

import (
	"context"
)

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	activeOutputs := make([]chan string, len(outputs))
	copy(activeOutputs, outputs)
	i := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if len(activeOutputs) == 0 {
				return nil
			}

			sent := false
			for attempts := 0; attempts < len(activeOutputs) && !sent; attempts++ {
				ch := activeOutputs[i]
				select {
				case ch <- data:
					sent = true
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				i = (i + 1) % len(activeOutputs)
			}
		}
	}
}
