package main

import (
	"fmt"
	"os"

	"task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Используем Print без \n, как в примерах
	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
