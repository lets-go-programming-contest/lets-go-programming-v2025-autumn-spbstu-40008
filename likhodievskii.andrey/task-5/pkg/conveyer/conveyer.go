package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var errChannelNotFound = errors.New("chan not found")

const undefinedValue = "undefined"

type Conveyer struct {
	size           int
	channelsByName map[string]chan string
	handlerList    []func(executionContext context.Context) error
	mutex          sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:           size,
		channelsByName: make(map[string]chan string),
		mutex:          sync.RWMutex{},
		handlerList:    []func(executionContext context.Context) error{},
	}
}

func (conveyer *Conveyer) register(channelName string) chan string {
	if channel, channelExists := conveyer.channelsByName[channelName]; channelExists {
		return channel
	}

	channel := make(chan string, conveyer.size)
	conveyer.channelsByName[channelName] = channel

	return channel
}

func (conveyer *Conveyer) RegisterDecorator(
	handlerFunction func(
		executionContext context.Context,
		inputChannel chan string,
		outputChannel chan string,
	) error,
	inputChannelName string,
	outputChannelName string,
) {
	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	outputChannel := conveyer.register(outputChannelName)
	inputChannel := conveyer.register(inputChannelName)

	conveyer.handlerList = append(conveyer.handlerList, func(context context.Context) error {
		return handlerFunction(context, inputChannel, outputChannel)
	})
}

func (conveyer *Conveyer) RegisterMultiplexer(
	handlerFunction func(
		executionContext context.Context,
		inputChannels []chan string,
		outputChannel chan string,
	) error,
	inputs []string,
	output string,
) {
	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = conveyer.register(name)
	}

	conveyer.handlerList = append(conveyer.handlerList, func(ctx context.Context) error {
		return handlerFunction(ctx, inChans, conveyer.register(output))
	})
}

func (conveyer *Conveyer) RegisterSeparator(
	callback func(
		context context.Context,
		inputChannel chan string,
		outputChannels []chan string,
	) error,
	input string,
	outputs []string,
) {
	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = conveyer.register(name)
	}

	inputChannel := conveyer.register(input)

	conveyer.handlerList = append(conveyer.handlerList, func(context context.Context) error {
		return callback(context, inputChannel, outChans)
	})
}

func (conveyer *Conveyer) Run(executionContext context.Context) error {
	defer func() {
		conveyer.mutex.RLock()
		defer conveyer.mutex.RUnlock()

		for _, channel := range conveyer.channelsByName {
			close(channel)
		}
	}()

	errorGroup, operationContext := errgroup.WithContext(executionContext)

	conveyer.mutex.RLock()

	for _, handler := range conveyer.handlerList {
		errorGroup.Go(func() error {
			return handler(operationContext)
		})
	}

	conveyer.mutex.RUnlock()

	if runError := errorGroup.Wait(); runError != nil {
		return fmt.Errorf("run pipeline: %w", runError)
	}

	return nil
}

func (conveyer *Conveyer) Send(inputChannelName string, data string) error {
	conveyer.mutex.RLock()

	channel, channelExists := conveyer.channelsByName[inputChannelName]

	conveyer.mutex.RUnlock()

	if !channelExists {
		return errChannelNotFound
	}

	channel <- data

	return nil
}

func (conveyer *Conveyer) Recv(outputChannelName string) (string, error) {
	conveyer.mutex.RLock()

	channel, channelExists := conveyer.channelsByName[outputChannelName]

	conveyer.mutex.RUnlock()

	if !channelExists {
		return "", errChannelNotFound
	}

	data, channelExists := <-channel
	if !channelExists {
		return undefinedValue, nil
	}

	return data, nil
}
