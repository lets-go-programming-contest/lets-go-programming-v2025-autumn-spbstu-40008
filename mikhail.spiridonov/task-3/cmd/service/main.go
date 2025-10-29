package main

import (
	"flag"
	"github.com/mordw1n/task-3/internal/config"
	"github.com/mordw1n/task-3/internal/parser"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config YAML file")
	flag.Parse()

	cfg := config.LoadFile(*configPath)
	parser.ParseAndSortXML(cfg.InputFile, cfg.OutputFile)
}
