package main

import (
	"flag"
	"fmt"

	"github.com/task-3/config"
	"github.com/task-3/internal/data"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "Path to the YAML configuration file")
	flag.Parse()

	if configPath == "" {
		panic("Flag --config must be set to the path of the YAML configuration file.")
	}

	cfg, err := config.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Configuration reading failed: %v", err))
	}

	if cfg.Input == "" || cfg.Output == "" {
		panic("Configuration file must contain 'input-file' and 'output-file' paths.")
	}

	xmlData, err := data.DecodeXML(cfg)
	if err != nil {
		panic(fmt.Sprintf("XML decoding failed: %v", err))
	}

	sortedCurrencies := data.ProcessAndSortCurrencies(xmlData)

	err = data.CreateAndWriteJSON(cfg.Output, sortedCurrencies)
	if err != nil {
		panic(fmt.Sprintf("JSON writing failed: %v", err))
	}

	fmt.Printf("Successfully processed %d currencies and saved results to '%s'\n", len(sortedCurrencies), cfg.Output)
}
