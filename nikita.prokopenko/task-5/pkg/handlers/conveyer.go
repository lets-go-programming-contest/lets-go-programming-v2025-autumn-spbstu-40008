package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Conveyer interface {
	RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string)
	RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyer struct {
	size      int
	pipes     map[string]chan string
	tasks     []taskDef
	lock      sync.RWMutex
}

type taskDef struct {
	kind     string
	work     interface{}
	inTags   []string
	outTags  []string
}

func New(size int) Conveyer {
	return &conveyer{
		size:  size,
		pipes: make(map[string]chan string),
		tasks: make([]taskDef, 0),
	}
}

func (c *conveyer) makePipe(tag string) chan string {
	c.lock.Lock()
	defer c.lock.Unlock()

	if pipe, ok := c.pipes[tag]; ok {
		return pipe
	}
	
	pipe := make(chan string, c.size)
	c.pipes[tag] = pipe
	return pipe
}

func (c *conveyer) RegisterDecorator(work func(ctx context.Context, input chan string, output chan string) error, inTag string, outTag string) {
	c.tasks = append(c.tasks, taskDef{
		kind:    "decorator",
		work:    work,
		inTags:  []string{inTag},
		outTags: []string{outTag},
	})
}

func (c *conveyer) RegisterMultiplexer(work func(ctx context.Context, inputs []chan string, output chan string) error, inTags []string, outTag string) {
	c.tasks = append(c.tasks, taskDef{
		kind:    "multiplexer",
		work:    work,
		inTags:  inTags,
		outTags: []string{outTag},
	})
}

func (c *conveyer) RegisterSeparator(work func(ctx context.Context, input chan string, outputs []chan string) error, inTag string, outTags []string) {
	c.tasks = append(c.tasks, taskDef{
		kind:    "separator",
		work:    work,
		inTags:  []string{inTag},
		outTags: outTags,
	})
}

func (c *conveyer) Run(ctx context.Context) error {
	for _, task := range c.tasks {
		for _, in := range task.inTags {
			c.makePipe(in)
		}
		for _, out := range task.outTags {
			c.makePipe(out)
		}
	}

	var workers sync.WaitGroup
	errorPipe := make(chan error, len(c.tasks))

	for _, task := range c.tasks {
		workers.Add(1)
		go func(t taskDef) {
			defer workers.Done()
			
			var runErr error
			switch t.kind {
			case "decorator":
				workFunc := t.work.(func(ctx context.Context, input chan string, output chan string) error)
				inPipe := c.makePipe(t.inTags[0])
				outPipe := c.makePipe(t.outTags[0])
				runErr = workFunc(ctx, inPipe, outPipe)
				close(outPipe)
				
			case "multiplexer":
				workFunc := t.work.(func(ctx context.Context, inputs []chan string, output chan string) error)
				inPipes := make([]chan string, len(t.inTags))
				for i, tag := range t.inTags {
					inPipes[i] = c.makePipe(tag)
				}
				outPipe := c.makePipe(t.outTags[0])
				runErr = workFunc(ctx, inPipes, outPipe)
				close(outPipe)
				
			case "separator":
				workFunc := t.work.(func(ctx context.Context, input chan string, outputs []chan string) error)
				inPipe := c.makePipe(t.inTags[0])
				outPipes := make([]chan string, len(t.outTags))
				for i, tag := range t.outTags {
					outPipes[i] = c.makePipe(tag)
				}
				runErr = workFunc(ctx, inPipe, outPipes)
				for _, outPipe := range outPipes {
					close(outPipe)
				}
			}
			
			if runErr != nil {
				select {
				case errorPipe <- runErr:
				default:
				}
			}
		}(task)
	}

	workers.Wait()
	close(errorPipe)

	for err := range errorPipe {
		return fmt.Errorf("task failure: %w", err)
	}

	return nil
}

func (c *conveyer) Send(pipeTag string, data string) error {
	c.lock.RLock()
	pipe, exists := c.pipes[pipeTag]
	c.lock.RUnlock()

	if !exists {
		return errors.New("chan not found")
	}

	select {
	case pipe <- data:
		return nil
	default:
		return errors.New("pipe full")
	}
}

func (c *conveyer) Recv(pipeTag string) (string, error) {
	c.lock.RLock()
	pipe, exists := c.pipes[pipeTag]
	c.lock.RUnlock()

	if !exists {
		return "", errors.New("chan not found")
	}

	value, ok := <-pipe
	if !ok {
		return "undefined", nil
	}
	
	return value, nil
}