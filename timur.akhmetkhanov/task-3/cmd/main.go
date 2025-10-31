package main

import (
	"flag"
	"log"

	"github.com/task-3/app"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration YAML file")
	flag.Parse()

	cfg, err := app.LoadConfig(*configPath)
	handleError(err)

	log.Printf("Входной файл: %s", cfg.InputFile)
	log.Printf("Выходной файл: %s", cfg.OutputFile)

	err = app.Run(cfg)
	handleError(err)
}
