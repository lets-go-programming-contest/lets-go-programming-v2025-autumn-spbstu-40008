package main

import (
	"fmt"

	"github.com/gagysun/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}