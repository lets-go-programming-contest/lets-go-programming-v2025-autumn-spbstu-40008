package main

import (
	"flag"

	"github.com/mordw1n/task-3/internal/config"
	"github.com/mordw1n/task-3/internal/parser"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "path to config YAML file")
	flag.Parse()

	config, err := config.LoadFile(*configPath)
	if err != nil {
		panic(err)
	}

	err = parser.ParseAndSortXML(config.InputFile, config.OutputFile)
	if err != nil {
		panic(err)
	}
}
