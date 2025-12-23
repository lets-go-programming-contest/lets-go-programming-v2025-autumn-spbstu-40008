package main

import (
	"fmt"
	"os"

	"github.com/mordw1n/task-8/internal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
