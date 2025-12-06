package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mu         sync.RWMutex
	pipeMap    map[string]chan string       
	processors []func(context.Context) error 
	bufferSize int
}

func New(size int) *Conveyer {
	return &Conveyer{
		pipeMap:    make(map[string]chan string),
		processors: make([]func(context.Context) error, 0),
		bufferSize: size,
	}
}

func (c *Conveyer) getOrInitChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.pipeMap[name]; exists {
		return ch
	}

	newCh := make(chan string, c.bufferSize)
	c.pipeMap[name] = newCh
	return newCh
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	inCh := c.getOrInitChannel(inputName)
	outCh := c.getOrInitChannel(outputName)

	wrappedTask := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}

	c.mu.Lock()
	c.processors = append(c.processors, wrappedTask)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	inputChans := make([]chan string, 0, len(inputNames))
	for _, name := range inputNames {
		inputChans = append(inputChans, c.getOrInitChannel(name))
	}
	
	outCh := c.getOrInitChannel(outputName)

	wrappedTask := func(ctx context.Context) error {
		return fn(ctx, inputChans, outCh)
	}

	c.mu.Lock()
	c.processors = append(c.processors, wrappedTask)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	inCh := c.getOrInitChannel(inputName)
	
	outputChans := make([]chan string, 0, len(outputNames))
	for _, name := range outputNames {
		outputChans = append(outputChans, c.getOrInitChannel(name))
	}

	wrappedTask := func(ctx context.Context) error {
		return fn(ctx, inCh, outputChans)
	}

	c.mu.Lock()
	c.processors = append(c.processors, wrappedTask)
	c.mu.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	for _, proc := range c.processors {
		wg.Add(1)
		currentProc := proc
		
		go func() {
			defer wg.Done()
			if err := currentProc(ctx); err != nil {
				select {
				case errChan <- err:
					cancel() 
				default:
				}
			}
		}()
	}

	wg.Wait()

	c.mu.Lock()
	for _, ch := range c.pipeMap {
		close(ch)
	}
	c.mu.Unlock()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (c *Conveyer) Send(name string, data string) error {
	c.mu.RLock()
	ch, ok := c.pipeMap[name]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mu.RLock()
	ch, ok := c.pipeMap[name]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, isOpen := <-ch
	if !isOpen {
		return "undefined", nil 
	}

	return val, nil
}