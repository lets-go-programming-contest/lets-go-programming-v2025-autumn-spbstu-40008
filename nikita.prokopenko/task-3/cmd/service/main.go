package main

import (
	"flag"
	"log"
	"sort"

	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-3/internal/config"
	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-3/internal/currency"
)

func main() {
	cfgPath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	if *cfgPath == "" {
		log.Fatal("configuration file path is required")
	}

	settings, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("configuration loading failed: %v", err)
	}

	items, err := currency.ParseCurrencyFile(settings.InputFile)
	if err != nil {
		log.Fatalf("currency data parsing failed: %v", err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})

	if err := currency.ExportToJSON(items, settings.OutputFile); err != nil {
		log.Fatalf("JSON export failed: %v", err)
	}
}