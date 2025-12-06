package conveyer

import (
    "context"
    "errors"
    "fmt"
    "sync"
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type Conveyer interface {
    RegisterDecorator(fn DecoratorFunc, input string, output string)
    RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string)
    RegisterSeparator(fn SeparatorFunc, input string, outputs []string)
    Run(ctx context.Context) error
    Send(input string, data string) error
    Recv(output string) (string, error)
}

type conveyerImpl struct {
    size      int
    channels  map[string]chan string
    processes []processInfo
    lock      sync.RWMutex
}

type processInfo struct {
    typ      string
    fn       interface{}
    inNames  []string
    outNames []string
}

func New(size int) Conveyer {
    return &conveyerImpl{
        size:      size,
        channels:  make(map[string]chan string),
        processes: make([]processInfo, 0),
    }
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
    c.lock.Lock()
    defer c.lock.Unlock()

    if ch, ok := c.channels[name]; ok {
        return ch
    }
    
    ch := make(chan string, c.size)
    c.channels[name] = ch
    return ch
}

func (c *conveyerImpl) RegisterDecorator(fn DecoratorFunc, input string, output string) {
    c.processes = append(c.processes, processInfo{
        typ:      "decorator",
        fn:       fn,
        inNames:  []string{input},
        outNames: []string{output},
    })
}

func (c *conveyerImpl) RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string) {
    c.processes = append(c.processes, processInfo{
        typ:      "multiplexer",
        fn:       fn,
        inNames:  inputs,
        outNames: []string{output},
    })
}

func (c *conveyerImpl) RegisterSeparator(fn SeparatorFunc, input string, outputs []string) {
    c.processes = append(c.processes, processInfo{
        typ:      "separator",
        fn:       fn,
        inNames:  []string{input},
        outNames: outputs,
    })
}

func (c *conveyerImpl) Run(ctx context.Context) error {
    for _, proc := range c.processes {
        for _, in := range proc.inNames {
            c.getOrCreateChannel(in)
        }
        for _, out := range proc.outNames {
            c.getOrCreateChannel(out)
        }
    }

    var wg sync.WaitGroup
    errChan := make(chan error, len(c.processes))

    for _, proc := range c.processes {
        wg.Add(1)
        go func(p processInfo) {
            defer wg.Done()
            
            var runErr error
            switch p.typ {
            case "decorator":
                fn := p.fn.(DecoratorFunc)
                inChan := c.getOrCreateChannel(p.inNames[0])
                outChan := c.getOrCreateChannel(p.outNames[0])
                runErr = fn(ctx, inChan, outChan)
                close(outChan)
                
            case "multiplexer":
                fn := p.fn.(MultiplexerFunc)
                inChans := make([]chan string, len(p.inNames))
                for i, name := range p.inNames {
                    inChans[i] = c.getOrCreateChannel(name)
                }
                outChan := c.getOrCreateChannel(p.outNames[0])
                runErr = fn(ctx, inChans, outChan)
                close(outChan)
                
            case "separator":
                fn := p.fn.(SeparatorFunc)
                inChan := c.getOrCreateChannel(p.inNames[0])
                outChans := make([]chan string, len(p.outNames))
                for i, name := range p.outNames {
                    outChans[i] = c.getOrCreateChannel(name)
                }
                runErr = fn(ctx, inChan, outChans)
                for _, outChan := range outChans {
                    close(outChan)
                }
            }
            
            if runErr != nil {
                select {
                case errChan <- runErr:
                default:
                }
            }
        }(proc)
    }

    wg.Wait()
    close(errChan)

    for err := range errChan {
        return fmt.Errorf("handler error: %w", err)
    }

    return nil
}

func (c *conveyerImpl) Send(input string, data string) error {
    c.lock.RLock()
    ch, exists := c.channels[input]
    c.lock.RUnlock()

    if !exists {
        return errors.New("chan not found")
    }

    select {
    case ch <- data:
        return nil
    default:
        return errors.New("channel buffer is full")
    }
}

func (c *conveyerImpl) Recv(output string) (string, error) {
    c.lock.RLock()
    ch, exists := c.channels[output]
    c.lock.RUnlock()

    if !exists {
        return "", errors.New("chan not found")
    }

    value, ok := <-ch
    if !ok {
        return "undefined", nil
    }
    
    return value, nil
}