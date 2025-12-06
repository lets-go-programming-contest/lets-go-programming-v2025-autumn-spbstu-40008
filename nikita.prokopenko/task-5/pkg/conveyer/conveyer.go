// pkg/conveyer/conveyer.go
package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrMissingChannel = errors.New("chan not found")

type Pipeline struct {
	mu       sync.RWMutex
	pipes    map[string]chan string
	workers  []func(ctx context.Context) error
	bufSize  int
}

func CreatePipeline(capacity int) *Pipeline {
	return &Pipeline{
		pipes:   make(map[string]chan string),
		workers: []func(ctx context.Context) error{},
		bufSize: capacity,
	}
}

func (p *Pipeline) obtainPipe(name string) chan string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if pipe, found := p.pipes[name]; found {
		return pipe
	}

	pipe := make(chan string, p.bufSize)
	p.pipes[name] = pipe
	return pipe
}

func (p *Pipeline) RegisterDecorator(
	decorator func(ctx context.Context, input, output chan string) error,
	inputName, outputName string,
) {
	inputPipe := p.obtainPipe(inputName)
	outputPipe := p.obtainPipe(outputName)

	p.mu.Lock()
	p.workers = append(p.workers, func(ctx context.Context) error {
		return decorator(ctx, inputPipe, outputPipe)
	})
	p.mu.Unlock()
}

func (p *Pipeline) RegisterMultiplexer(
	merger func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	inputPipes := make([]chan string, len(inputNames))
	for i, name := range inputNames {
		inputPipes[i] = p.obtainPipe(name)
	}

	outputPipe := p.obtainPipe(outputName)

	p.mu.Lock()
	p.workers = append(p.workers, func(ctx context.Context) error {
		return merger(ctx, inputPipes, outputPipe)
	})
	p.mu.Unlock()
}

func (p *Pipeline) RegisterSeparator(
	distributor func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputPipe := p.obtainPipe(inputName)
	outputPipes := make([]chan string, len(outputNames))
	for i, name := range outputNames {
		outputPipes[i] = p.obtainPipe(name)
	}

	p.mu.Lock()
	p.workers = append(p.workers, func(ctx context.Context) error {
		return distributor(ctx, inputPipe, outputPipes)
	})
	p.mu.Unlock()
}

func (p *Pipeline) Send(pipeName string, content string) error {
	p.mu.RLock()
	pipe, exists := p.pipes[pipeName]
	p.mu.RUnlock()

	if !exists {
		return ErrMissingChannel
	}

	pipe <- content
	return nil
}

func (p *Pipeline) Receive(pipeName string) (string, error) {
	p.mu.RLock()
	pipe, exists := p.pipes[pipeName]
	p.mu.RUnlock()

	if !exists {
		return "", ErrMissingChannel
	}

	data, isOpen := <-pipe
	if !isOpen {
		return "undefined", nil
	}

	return data, nil
}

func (p *Pipeline) Execute(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	for _, worker := range p.workers {
		wg.Add(1)
		go func(task func(ctx context.Context) error) {
			defer wg.Done()
			if err := task(ctx); err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(worker)
	}

	wg.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()
	for _, pipe := range p.pipes {
		close(pipe)
	}

	return firstErr
}