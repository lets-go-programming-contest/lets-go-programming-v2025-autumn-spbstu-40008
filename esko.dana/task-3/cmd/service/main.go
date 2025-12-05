package main

import (
	"flag"
	"fmt"

	"esko.dana/task-3/iternal/config"
	"esko.dana/task-3/iternal/currency"
	"esko.dana/task-3/iternal/json"
	"esko.dana/task-3/iternal/xml"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "Path to configuration YAML file")
	flag.Parse()

	if configPath == "" {
		panic("Configuration file path is required. Use -config <path>")
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}

	valutes, err := xml.Parse(cfg.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Error parsing XML: %v", err))
	}

	sortedCurrencies, err := currency.ProcessAndSort(valutes)
	if err != nil {
		panic(fmt.Sprintf("Error processing/sorting currencies: %v", err))
	}

	err = json.Save(sortedCurrencies, cfg.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("Error saving JSON: %v", err))
	}

	fmt.Printf("Success! Processed %d currencies. Result saved to: %s\n",
		len(sortedCurrencies), cfg.OutputFile)
}
