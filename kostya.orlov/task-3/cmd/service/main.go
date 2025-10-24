package main

import (
	"flag"

	"github.com/TWChromeTW/task-3/config"
	"github.com/TWChromeTW/task-3/currency"
)

func main() {
	configFile := flag.String("config", "config/config.yaml", "config/config.yaml")
	outputFormat := flag.String("output-format", "json", "Output format: json, yaml, or xml")

	flag.Parse()

	cfg, err := config.Load(*configFile)

	if err != nil {
		panic(err)
	}

	valutes, err := currency.DecodeXML(cfg.InputFile)

	if err != nil {
		panic(err)
	}

	currency.SortValutes(valutes)

	err = currency.EncodeFile(valutes, *outputFormat, cfg.OutputFile)

	if err != nil {
		panic(err)
	}
}
