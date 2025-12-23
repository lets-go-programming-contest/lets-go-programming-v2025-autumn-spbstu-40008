package main

import (
	"fmt"

	"github.com/JDH-LR-994/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("load config error: %w", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
