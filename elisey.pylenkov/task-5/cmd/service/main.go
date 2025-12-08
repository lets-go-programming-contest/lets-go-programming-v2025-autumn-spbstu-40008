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

	err := conv.Send("nonexistent", "test")
	fmt.Printf("Send to nonexistent channel: %v\n", err)

	_, err = conv.Recv("nonexistent")
	fmt.Printf("Recv from nonexistent channel: %v\n", err)

	conv2 := conveyer.New(10)
	conv2.RegisterDecorator(handlers.PrefixDecoratorFunc, "input1", "decorated1")
	conv2.RegisterSeparator(handlers.SeparatorFunc, "decorated1", []string{"out1", "out2"})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go func() {
		if err := conv2.Run(ctx); err != nil {
			fmt.Println("Conveyer error:", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	messages := []string{"data1", "data2", "data3", "data4"}
	for _, v := range messages {
		if err := conv2.Send("input1", v); err != nil {
			fmt.Println("Send error:", err)
		}
	}

	for i := 0; i < 2; i++ {
		data, err := conv2.Recv("out1")
		if err != nil {
			fmt.Println("Recv out1 error:", err)
		} else {
			fmt.Println("out1:", data)
		}
	}

	for i := 0; i < 2; i++ {
		data, err := conv2.Recv("out2")
		if err != nil {
			fmt.Println("Recv out2 error:", err)
		} else {
			fmt.Println("out2:", data)
		}
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Program finished")
}
