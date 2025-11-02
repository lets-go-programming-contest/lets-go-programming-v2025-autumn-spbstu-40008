package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/rabbitdfs/task-3/internal/config"
	"github.com/rabbitdfs/task-3/internal/currency"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "yaml config path")
	flag.Parse()

	if configPath == "" {
		panic("no config provided")
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	data, err := currency.ReadXML(cfg.InputFile)
	if err != nil {
		panic(fmt.Errorf("read xml: %w", err))
	}

	sort.Slice(data.Currencies, func(i, j int) bool {
		return float64(data.Currencies[i].Value) > float64(data.Currencies[j].Value)
	})

	if err := currency.WriteJSON(data.Currencies, cfg.OutputFile); err != nil {
		panic(fmt.Errorf("write json: %w", err))
	}
}
