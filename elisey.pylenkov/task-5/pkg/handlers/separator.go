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
				case <-ctx.Done():
					return ctx.Err()
				default:
					func() {
						defer func() {
							if r := recover(); r != nil {
								activeOutputs = append(activeOutputs[:i], activeOutputs[i+1:]...)
								if i >= len(activeOutputs) {
									i = 0
								}
							}
						}()
						select {
						case ch <- data:
							sent = true
						default:
						}
					}()
				}
				if !sent {
					i = (i + 1) % len(activeOutputs)
				}
			}
			if sent {
				i = (i + 1) % len(activeOutputs)
			}
		}
	}
}
