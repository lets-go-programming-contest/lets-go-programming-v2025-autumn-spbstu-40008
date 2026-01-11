package main

import (
	"fmt"
	"os"

	"github.com/TWChromeTW/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
