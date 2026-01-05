package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"esko.dana/task-3/internal/config"
	"esko.dana/task-3/internal/currency"
	"esko.dana/task-3/internal/json"
	"esko.dana/task-3/internal/xml"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "Path to configuration YAML file")
	flag.Parse()

	if configPath == "" {
		log.Fatal("Configuration file path is required. Use -config <path>")
	}

	// Проверяем существование файла конфига
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if cfg.InputPath == "" {
		log.Fatal("Input path is not specified in config")
	}

	valutes, err := xml.Parse(cfg.InputPath)
	if err != nil {
		log.Fatalf("Error parsing XML: %v", err)
	}

	sortedCurrencies, err := currency.ProcessAndSort(valutes)
	if err != nil {
		log.Fatalf("Error processing/sorting currencies: %v", err)
	}

	outputDir := filepath.Dir(cfg.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	err = json.Write(sortedCurrencies, cfg.OutputPath)
	if err != nil {
		log.Fatalf("Error saving JSON: %v", err)
	}

	fmt.Printf("Success! Processed %d currencies. Result saved to: %s\n",
		len(sortedCurrencies), cfg.OutputPath)
}
