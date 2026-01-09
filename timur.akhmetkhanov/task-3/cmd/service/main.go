package main

import (
	"flag"

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

	err = app.Run(cfg)
	handleError(err)
}
