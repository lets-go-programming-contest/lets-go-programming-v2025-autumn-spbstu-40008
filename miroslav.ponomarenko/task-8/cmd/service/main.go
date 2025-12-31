package main

import (
	"fmt"
	"log"

	"miroslav.ponomarenko/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
