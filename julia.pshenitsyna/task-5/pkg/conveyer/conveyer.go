package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channels map[string]chan string
	myMutex  sync.RWMutex
	size     int
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		myMutex:  sync.RWMutex{},
		size:     size,
		handlers: make([]func(ctx context.Context) error, 0),
	}
}

func (conv *Conveyer) createChan(name string) chan string {
	conv.myMutex.Lock()
	defer conv.myMutex.Unlock()

	if channel, okey := conv.channels[name]; okey {
		return channel
	}

	channel := make(chan string, conv.size)
	conv.channels[name] = channel

	return channel
}

func (conv *Conveyer) RegisterDecorator(handlerFn func(
	ctx context.Context,
	input chan string,
	output chan string,
) error,
	input string,

	output string,
) {
	inputChannel := conv.createChan(input)
	outputChannel := conv.createChan(output)

	conv.myMutex.Lock()
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputChannel, outputChannel)
	})
	conv.myMutex.Unlock()
}

func (conv *Conveyer) RegisterMultiplexer(handlerFn func(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))

	for i := range inputs {
		inputChannels[i] = conv.createChan(inputs[i])
	}

	outputChannel := conv.createChan(output)

	conv.myMutex.Lock()
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputChannels, outputChannel)
	})
	conv.myMutex.Unlock()
}

func (conv *Conveyer) RegisterSeparator(handlerFn func(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error,
	input string,
	outputs []string,
) {
	outputChannels := make([]chan string, len(outputs))

	for i := range outputs {
		outputChannels[i] = conv.createChan(outputs[i])
	}

	inputChannel := conv.createChan(input)

	conv.myMutex.Lock()
	defer conv.myMutex.Unlock()

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputChannel, outputChannels)
	})
}

func (conv *Conveyer) Run(ctx context.Context) error {
	conv.myMutex.RLock()
	handlers := make([]func(context.Context) error, len(conv.handlers))
	copy(handlers, conv.handlers)
	conv.myMutex.RUnlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup

	errCh := make(chan error, 1)

	for _, handler := range handlers {
		waitGroup.Add(1)

		go func(handlerFunc func(context.Context) error) {
			defer waitGroup.Done()

			if err := handlerFunc(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				cancel()
			}
		}(handler)
	}

	done := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		<-done
		conv.closeAll()

		return err
	case <-done:
		conv.closeAll()

		return nil
	case <-ctx.Done():
		<-done
		conv.closeAll()

		return nil
	}
}

func (conv *Conveyer) Send(input string, data string) error {
	conv.myMutex.RLock()
	channel, okey := conv.channels[input]
	conv.myMutex.RUnlock()

	if !okey {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (conv *Conveyer) Recv(output string) (string, error) {
	conv.myMutex.RLock()
	channel, okey := conv.channels[output]
	conv.myMutex.RUnlock()

	if !okey {
		return "", ErrChanNotFound
	}

	val, okey := <-channel
	if !okey {
		return "undefined", nil
	}

	return val, nil
}

func (conv *Conveyer) closeAll() {
	conv.myMutex.Lock()
	defer conv.myMutex.Unlock()

	for name, channel := range conv.channels {
		if channel != nil {
			close(channel)
			conv.channels[name] = nil
		}
	}
}
