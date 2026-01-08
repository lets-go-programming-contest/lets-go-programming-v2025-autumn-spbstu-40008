package main

import (
	"flag"
	"log"
	"sort"

	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-3/internal/config"
	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-3/internal/currency"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "configs/config.yaml", "path to config")
	flag.Parse()

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	items, err := currency.DecodeXMLFile(cfg.InputFile)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].NumCode > items[j].NumCode
	})

	if err := currency.SaveAsJSON(cfg.OutputFile, items); err != nil {
		log.Fatal(err)
	}
}
