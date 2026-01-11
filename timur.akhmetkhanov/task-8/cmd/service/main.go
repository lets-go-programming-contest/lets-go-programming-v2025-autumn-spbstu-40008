package main

import (
	"fmt"
	"log"
	"task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
