package main

import (
	"fmt"
	"log"

	"github.com/Ilya-Er0fick/task-8/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	
	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
