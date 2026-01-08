package main

import (
	"flag"
	"log"

	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-3/internal/config"
	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-3/internal/currency"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "configs/config.yaml", "")
	flag.Parse()

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	items, err := currency.DecodeXMLFile(cfg.InputFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := currency.SaveAsJSON(cfg.OutputFile, items); err != nil {
		log.Fatal(err)
	}
}
