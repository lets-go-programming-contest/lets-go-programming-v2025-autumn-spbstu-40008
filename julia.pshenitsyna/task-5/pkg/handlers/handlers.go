package handlers

import (
	"context"
	"errors"
	"strings"
)

func PrefixDecoratorFunc (ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(val, "decorated:") {
				val = "decorated:" + val
			}

			output <- val
		}
	}
}

func SeparatorFunc (ctx context.Context, input chan string, outputs []chan string) error {
	cnt := 0

	for {
		select{
		case <- ctx.Done():
			return nil

		case val, ok := <- input:
			if !ok{

				return nil
			}

			outputs[cnt] <- val
			cnt++

			if cnt >= len(outputs){
				cnt = 0
			}
		}
	}
}

func MultiplexerFunc (ctx context.Context, inputs []chan string, output chan string) error {
	for i := 0; i < len(inputs); i++{
		in := inputs[i]

		go func(channel chan string){
			for {
				select{
				case <- ctx.Done():
					return

				case val, ok := <- channel:
					if !ok{
						return
					}

					if strings.Contains(val, "no multiplexer"){
						continue
					}

					output <- val
				}
			}
		} (in)
	}
	<- ctx.Done()

	return nil
}
