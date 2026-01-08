package main

import (
	"flag"
	"log"
	"os"
	"sort"

	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-3/internal/appconfig"
	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-3/internal/conversion"
)

func main() {
	cfgPath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	if *cfgPath == "" {
		log.Fatal("configuration file path is required")
	}

	settings, err := appconfig.New(*cfgPath)
	if err != nil {
		log.Fatalf("configuration loading failed: %v", err)
	}

	rates, err := conversion.ParseCurrencyData(settings.SourceFile)
	if err != nil {
		log.Fatalf("currency data parsing failed: %v", err)
	}

	currencies := rates.GetAllCurrencies()
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	if err := conversion.ExportToJSON(currencies, settings.TargetFile); err != nil {
		log.Fatalf("JSON export failed: %v", err)
	}
}