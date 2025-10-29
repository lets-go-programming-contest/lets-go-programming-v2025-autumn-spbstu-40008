package main

import (
	"flag"
	"os"
	"github.com/mordw1n/task-3/config"
	"github.com/mordw1n/task-3/parser"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config YAML file")
	flag.Parse()
	
	cfg := config.LoadFile(*configPath)
	parser.ParseAndSortXML(cfg.InputFile, cfg.OutputFile)
}
