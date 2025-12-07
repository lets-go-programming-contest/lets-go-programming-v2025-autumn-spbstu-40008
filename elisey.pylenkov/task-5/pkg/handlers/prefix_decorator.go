package handlers

import "context"

func PrefixDecoratorFunc(prefix string) DecoratorFunc[string] {
	return func(ctx context.Context, s string) (string, error) {
		return prefix + s, nil
	}
}
