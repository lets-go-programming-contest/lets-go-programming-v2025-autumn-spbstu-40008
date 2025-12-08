package main

import (
	"context"
	"fmt"
	"time"

	"task-5/pkg/conveyer"
	"task-5/pkg/handlers"

	_ "github.com/stretchr/testify/require"
	_ "golang.org/x/sync/errgroup"
)

func main() {
	conv := conveyer.New(10)

	conv.RegisterDecorator(handlers.PrefixDecoratorFunc, "input1", "decorated1")
	conv.RegisterSeparator(handlers.SeparatorFunc, "decorated1", []string{"out1", "out2"})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := conv.Run(ctx); err != nil {
			fmt.Println("Conveyer error:", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	messages := []string{"data1", "data2", "data3", "data4"}
	for _, v := range messages {
		if err := conv.Send("input1", v); err != nil {
			fmt.Println("Send error:", err)
		}
	}

	for i := 0; i < 2; i++ {
		data, err := conv.Recv("out1")
		if err != nil {
			fmt.Println("Recv out1 error:", err)
		} else {
			fmt.Println("out1:", data)
		}
	}

	for i := 0; i < 2; i++ {
		data, err := conv.Recv("out2")
		if err != nil {
			fmt.Println("Recv out2 error:", err)
		} else {
			fmt.Println("out2:", data)
		}
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Program finished")
}
