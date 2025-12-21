package main

import (
	"fmt"
	"log"

	"github.com/AliseMarfina/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
