package handlers

import (
    "context"
    "errors"
    "strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
    const prefix = "decorated: "
    const errorSubstring = "no decorator"
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case data, ok := <-input:
            if !ok {
                return nil
            }
            
            if strings.Contains(data, errorSubstring) {
                return errors.New("can't be decorated")
            }
            
            if !strings.HasPrefix(data, prefix) {
                data = prefix + data
            }
            
            select {
            case output <- data:
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }
}


func SeparatorFunc() error {
}

func MultiplexerFunc() error {
}