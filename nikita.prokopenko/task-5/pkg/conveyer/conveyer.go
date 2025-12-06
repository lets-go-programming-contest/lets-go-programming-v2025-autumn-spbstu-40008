// pkg/conveyer/conveyer.go
package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChannelMissing = errors.New("chan not found")

type Flow struct {
	lock     sync.RWMutex
	streams  map[string]chan string
	routines []func(ctx context.Context) error
	capacity int
}

func CreateFlow(buffer int) *Flow {
	return &Flow{
		streams:  make(map[string]chan string),
		routines: make([]func(ctx context.Context) error, 0),
		capacity: buffer,
	}
}

func (f *Flow) accessChannel(name string) chan string {
	f.lock.Lock()
	defer f.lock.Unlock()

	if stream, found := f.streams[name]; found {
		return stream
	}

	stream := make(chan string, f.capacity)
	f.streams[name] = stream
	return stream
}

func (f *Flow) findChannel(name string) (chan string, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	stream, found := f.streams[name]
	return stream, found
}

func (f *Flow) RegisterDecorator(
	transform func(ctx context.Context, source, target chan string) error,
	sourceName string,
	targetName string,
) {
	sourceCh := f.accessChannel(sourceName)
	targetCh := f.accessChannel(targetName)

	f.lock.Lock()
	f.routines = append(f.routines, func(ctx context.Context) error {
		return transform(ctx, sourceCh, targetCh)
	})
	f.lock.Unlock()
}

func (f *Flow) RegisterMultiplexer(
	combine func(ctx context.Context, sources []chan string, target chan string) error,
	sourceNames []string,
	targetName string,
) {
	sources := make([]chan string, len(sourceNames))
	for i, name := range sourceNames {
		sources[i] = f.accessChannel(name)
	}

	targetCh := f.accessChannel(targetName)

	f.lock.Lock()
	f.routines = append(f.routines, func(ctx context.Context) error {
		return combine(ctx, sources, targetCh)
	})
	f.lock.Unlock()
}

func (f *Flow) RegisterSeparator(
	split func(ctx context.Context, source chan string, targets []chan string) error,
	sourceName string,
	targetNames []string,
) {
	sourceCh := f.accessChannel(sourceName)
	targets := make([]chan string, len(targetNames))
	for i, name := range targetNames {
		targets[i] = f.accessChannel(name)
	}

	f.lock.Lock()
	f.routines = append(f.routines, func(ctx context.Context) error {
		return split(ctx, sourceCh, targets)
	})
	f.lock.Unlock()
}

func (f *Flow) Send(name string, data string) error {
	stream, found := f.findChannel(name)
	if !found {
		return ErrChannelMissing
	}

	stream <- data
	return nil
}

func (f *Flow) Receive(name string) (string, error) {
	stream, found := f.findChannel(name)
	if !found {
		return "", ErrChannelMissing
	}

	message, isAlive := <-stream
	if !isAlive {
		return "undefined", nil
	}

	return message, nil
}

func (f *Flow) Execute(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var group sync.WaitGroup
	var firstErr error
	var once sync.Once

	for _, routine := range f.routines {
		group.Add(1)
		go func(task func(ctx context.Context) error) {
			defer group.Done()
			if err := task(ctx); err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
			}
		}(routine)
	}

	group.Wait()

	f.lock.Lock()
	defer f.lock.Unlock()
	for _, stream := range f.streams {
		close(stream)
	}

	return firstErr
}