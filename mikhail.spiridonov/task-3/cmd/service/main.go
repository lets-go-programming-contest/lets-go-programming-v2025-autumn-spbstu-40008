package main

import (
	"flag"
	"github.com/mordw1n/task-3/config"
	"github.com/mordw1n/task-3/parser"
)

func main() {
	inputFile := flag.String("input", "in/currencies.xml", "path to input XML file")
	outputFile := flag.String("output", "out/converted.json", "path to output JSON file")
	flag.Parse()
	
	parser.ParseAndSortXML(*inputFile, *outputFile)
}
