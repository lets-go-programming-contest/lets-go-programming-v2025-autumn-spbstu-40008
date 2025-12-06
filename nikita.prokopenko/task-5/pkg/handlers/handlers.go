package handlers

import (
	"context"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case value, open := <-input:
			if !open {
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
	if len(outputs) == 0 {
		return nil
	}
	
	pos := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case value, open := <-input:
			if !open {
				return nil
			}
			
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case outputs[pos] <- value:
					pos = (pos + 1) % len(outputs)
					break
				}
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	failChan := make(chan error, len(inputs))
	
	for _, in := range inputs {
		go func(inputPipe chan string) {
			for {
				select {
				case <-ctx.Done():
					failChan <- ctx.Err()
					return
				case value, open := <-inputPipe:
					if !open {
						failChan <- nil
						return
					}
					
					if strings.Contains(value, "no multiplexer") {
						continue
					}
					
					select {
					case <-ctx.Done():
						failChan <- ctx.Err()
						return
					case output <- value:
					}
				}
			}
		}(in)
	}
	
	for i := 0; i < len(inputs); i++ {
		if err := <-failChan; err != nil {
			return err
		}
	}
	
	return nil
}

type DecorateFail string

func (e DecorateFail) Error() string {
	return string(e)
}