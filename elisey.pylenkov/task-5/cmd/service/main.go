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

	conv.RegisterDecorator(
		handlers.PrefixDecoratorFunc,
		"input1",
		"decorated1",
	)

	conv.RegisterSeparator(
		handlers.SeparatorFunc,
		"decorated1",
		[]string{"out1", "out2"},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := conv.Run(ctx); err != nil && err != context.Canceled {
			fmt.Printf("Conveyer error: %v\n", err)
		}
		fmt.Println("Conveyer stopped")
	}()

	time.Sleep(50 * time.Millisecond)

	go func() {
		conv.Send("input1", "test data 1")
		conv.Send("input1", "test data 2")
		conv.Send("input1", "test data 3")
		fmt.Println("All data sent")
	}()

	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 2; i++ {
			data, err := conv.Recv("out1")
			if err != nil {
				fmt.Printf("Error reading from out1: %v\n", err)
			} else {
				fmt.Printf("out1: %s\n", data)
			}
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1; i++ {
			data, err := conv.Recv("out2")
			if err != nil {
				fmt.Printf("Error reading from out2: %v\n", err)
			} else {
				fmt.Printf("out2: %s\n", data)
			}
		}
		done <- true
	}()

	<-done
	<-done

	fmt.Println("All data received")
	cancel()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Program finished")
}
