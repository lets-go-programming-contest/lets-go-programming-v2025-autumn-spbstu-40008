package main

import (
	"flag"
	"github.com/mordw1n/task-3/config"
	"github.com/mordw1n/task-3/parser"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()
	
	cfg := config.LoadFile(*configPath)
	parser.ParseAndSortXML(cfg.InputFile, cfg.OutputFile)
}
