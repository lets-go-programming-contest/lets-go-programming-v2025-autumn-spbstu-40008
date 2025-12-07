package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct { 
	channels map[string] chan string
	myMutex sync.RWMutex
	size int
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer { 
	return &Conveyer{
		channels: make(map[string]chan string),
		myMutex: sync.RWMutex{},
		size: size,
		handlers: make([]func(ctx context.Context) error, 0),
	}
}

func (conv *Conveyer) createChan(name string) chan string{
	if channel, ok := conv.channels[name]; ok {
		return channel
	}

	channel := make(chan string, conv.size)
	conv.channels[name] = channel

	return channel
}

/*func (conv *Conveyer) getChan(name string) chan string{

}*/

func (conv *Conveyer)RegisterDecorator(fn func(
			ctx context.Context,
			input chan string,
			output chan string,
		) error, 
		input string,

		output string,
	){
		inputChannel := conv.createChan(input)
		outputChannel := conv.createChan(output)
		
		conv.myMutex.Lock()
		defer conv.myMutex.Unlock()

		conv.handlers = append(conv.handlers, func(ctx context.Context) error {
			return fn(ctx, inputChannel, outputChannel)
		})
	}

func (conv *Conveyer)RegisterMultiplexer(fn func(
			ctx context.Context,
			inputs []chan string,
			output chan string,
		) error, 
		inputs []string,
		output string,
	){
		conv.myMutex.Lock()
		inputChannels := make([]chan string, len(inputs))
		for i:= 0; i<len(inputs); i++{
			inputChannels[i] = conv.createChan(inputs[i])
		}
		outputChannel := conv.createChan(output)
		defer conv.myMutex.Unlock()
		
		conv.myMutex.Lock()
		defer conv.myMutex.Unlock()

		conv.handlers = append(conv.handlers, func(ctx context.Context) error {
			return fn(ctx, inputChannels, outputChannel)
		})
	}

func (conv *Conveyer)RegisterSeparator(fn func(
			ctx context.Context,
			input chan string,
			outputs []chan string,
		) error, 
		input string,
		outputs []string,
	){
		conv.myMutex.Lock()
		outputChannels := make([]chan string, len(outputs))
		for i:= 0; i<len(outputs); i++{
			outputChannels[i] = conv.createChan(outputs[i])
		}
		inputChannel := conv.createChan(input)
		defer conv.myMutex.Unlock()

		conv.handlers = append(conv.handlers, func(ctx context.Context) error {
			return fn(ctx, inputChannel, outputChannels)
		})
	}
	
func (conv *Conveyer) Run(ctx context.Context) error {
	if len(conv.handlers) == 0{
		return nil
	}

	var waitGroup sync.WaitGroup
	errChan := make(chan error, 1)

	for _, h := range conv.handlers{
		waitGroup.Add(1)
		go func (handler func(context.Context) error)  {
			defer waitGroup.Done()
			if err := handler(ctx); err != nil{
				select{
				case errChan <- err:
				default:

				}
			}
		} (h)
	}

	go func ()  {
		waitGroup.Wait()
		close(errChan)
	}()

	select{
	case err := <- errChan:
		if err != nil{
			conv.closeAll()
			return err
		}
		conv.closeAll()
	case <- ctx.Done():
		waitGroup.Wait()
		conv.closeAll()
		return ctx.Err()
	}
	return nil
}

func (conv *Conveyer) Send(input string, data string) error {
	conv.myMutex.RLock()
	channel, ok := conv.channels[input]
	conv.myMutex.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (conv *Conveyer) Recv(output string) (string, error) {
	conv.myMutex.RLock()
	channel, ok := conv.channels[output]
	conv.myMutex.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, ok := <- channel
	if !ok {
		return "undefined", nil
	}

	return val, nil
}

func (conv *Conveyer) closeAll (){
	conv.myMutex.Lock()
	defer conv.myMutex.Unlock()

	for _, channel := range conv.channels{
		close(channel)
	}
}
