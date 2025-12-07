package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoDecorator  = errors.New("can't be decorated")
	ErrEmptyOutputs = errors.New("empty outputs")
)

const (
	noDecorator     = "no decorator"
	noMultiplexer   = "no multiplexer"
	decoratedPrefix = "decorated: "
)

func PrefixDecoratorFunc(
	executionContext context.Context,
	inputChannel chan string,
	outputChannel chan string,
) error {
	for {
		select {
		case currentString, channelIsOpen := <-inputChannel:
			if !channelIsOpen {
				return nil
			}

			if strings.Contains(currentString, noDecorator) {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(currentString, decoratedPrefix) {
				currentString = decoratedPrefix + currentString
			}

			select {
			case outputChannel <- currentString:

			case <-executionContext.Done():
				return nil
			}

		case <-executionContext.Done():
			return nil
		}
	}
}

func SeparatorFunc(
	executionContext context.Context,
	inputChannel chan string,
	outputChannels []chan string,
) error {
	if len(outputChannels) == 0 {
		return ErrEmptyOutputs
	}

	currentIndex := 0

	for {
		select {
		case currentString, channelIsOpen := <-inputChannel:
			if !channelIsOpen {
				return nil
			}

			select {
			case outputChannels[currentIndex] <- currentString:

			case <-executionContext.Done():
				return nil
			}

			currentIndex = (currentIndex + 1) % len(outputChannels)

		case <-executionContext.Done():
			return nil
		}
	}
}

func MultiplexerFunc(
	executionContext context.Context,
	inputChannels []chan string,
	outputChannel chan string,
) error {
	waitForAllInputs := sync.WaitGroup{}

	waitForAllInputs.Add(len(inputChannels))

	for _, currentChannel := range inputChannels {
		go func() {
			defer waitForAllInputs.Done()

			for {
				select {
				case currentString, channelIsOpen := <-currentChannel:
					if !channelIsOpen {
						return
					}

					if strings.Contains(currentString, noMultiplexer) {
						continue
					}

					select {
					case outputChannel <- currentString:

					case <-executionContext.Done():
						return
					}

				case <-executionContext.Done():
					return
				}
			}
		}()
	}

	waitForAllInputs.Wait()

	return nil
}
