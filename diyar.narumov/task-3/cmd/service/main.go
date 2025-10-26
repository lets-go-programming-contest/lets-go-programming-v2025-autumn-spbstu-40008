package main

import (
	"flag"

	"github.com/narumov-diyar/task-3/internal/app"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	outputFormat := flag.String("output-format", "json", "output format: json, yaml, xml")
	flag.Parse()

	if *configPath == "" {
		panic("--config flag is required")
	}

	if err := app.Run(*configPath, *outputFormat); err != nil {
		panic(err)
	}
}
