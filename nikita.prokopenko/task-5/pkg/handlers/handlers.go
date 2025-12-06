package handlers

import (
    "context"
    "strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
    defer close(output)
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case value, ok := <-input:
            if !ok {
                return nil
            }
            
            if strings.Contains(value, "no decorator") {
                return DecorateFail("can't be decorated")
            }
            
            if !strings.HasPrefix(value, "decorated: ") {
                value = "decorated: " + value
            }
            
            select {
            case <-ctx.Done():
                return ctx.Err()
            case output <- value:
            }
        }
    }
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
    defer func() {
        for _, out := range outputs {
            close(out)
        }
    }()
    
    if len(outputs) == 0 {
        return nil
    }
    
    index := 0
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case value, ok := <-input:
            if !ok {
                return nil
            }
            
            for {
                select {
                case <-ctx.Done():
                    return ctx.Err()
                case outputs[index] <- value:
                    index = (index + 1) % len(outputs)
                    break
                }
            }
        }
    }
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
    defer close(output)
    
    if len(inputs) == 0 {
        return nil
    }
    
    errChan := make(chan error, len(inputs))
    
    for _, in := range inputs {
        go func(inputChan chan string) {
            for {
                select {
                case <-ctx.Done():
                    errChan <- ctx.Err()
                    return
                case value, ok := <-inputChan:
                    if !ok {
                        errChan <- nil
                        return
                    }
                    
                    if strings.Contains(value, "no multiplexer") {
                        continue
                    }
                    
                    select {
                    case <-ctx.Done():
                        errChan <- ctx.Err()
                        return
                    case output <- value:
                    }
                }
            }
        }(in)
    }
    
    for i := 0; i < len(inputs); i++ {
        if err := <-errChan; err != nil {
            return err
        }
    }
    
    return nil
}

type DecorateFail string

func (e DecorateFail) Error() string {
    return string(e)
}