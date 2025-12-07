package handlers

import (
	"context"
	"time"
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

			attempts := 0
			sent := false
			for attempts < len(outputs) && !sent {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case outputs[i] <- data:
					i = (i + 1) % len(outputs)
					sent = true
				case <-time.After(10 * time.Millisecond):
					i = (i + 1) % len(outputs)
					attempts++
				}
			}
		case <-time.After(100 * time.Millisecond):
			continue
		}
	}
}
