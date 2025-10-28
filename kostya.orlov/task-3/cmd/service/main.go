package main

import (
	"flag"

	"github.com/TWChromeTW/task-3/internal"
)

func main() {
	configFile := flag.String("config", "examples/config.yaml", "examples/config.yaml")
	outputFormat := flag.String("output-format", "json", "Output format: json, yaml, or xml")

	flag.Parse()

	cfg, err := internal.Load(*configFile)
	if err != nil {
		panic(err)
	}

	valutes, err := internal.DecodeXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	internal.SortValutes(valutes)

	err = internal.EncodeFile(valutes, *outputFormat, cfg.OutputFile)
	if err != nil {
		panic(err)
	}
}
