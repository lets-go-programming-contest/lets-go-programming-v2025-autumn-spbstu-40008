package main

import (
	"fmt"

	"github.com/kuzminykh.ulyana/task-8/internal/config"
)

func main() {
	if config.CurrentErr != nil {
		fmt.Printf("Config error: %v\n", config.CurrentErr)
		return
	}

	fmt.Println(config.Current.Environment, config.Current.LogLevel)
}
