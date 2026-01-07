package main

import (
	"encoding/json"
	"flag"
	"log"
	"sort"

	"github.com/task-3/internal/config"
	"github.com/task-3/internal/decoder"
	"github.com/task-3/internal/output"
	"github.com/task-3/internal/structures"
)

func sortByValue(valutes []structures.Valute) {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "Path to YAML configuration file")
	flag.Parse()

	if configPath == "" {
		log.Fatal("Config file path is required")
	}

	cfg, err := config.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	valCurs, err := decoder.DecodeXML(cfg)
	if err != nil {
		log.Fatalf("Failed to decode XML: %v", err)
	}

	sortByValue(valCurs.Valute)

	jsonData, err := json.MarshalIndent(valCurs.Valute, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	outputFile, err := output.CreateFile(cfg.Output)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputFile.Close()

	if _, err := outputFile.Write(jsonData); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}

