package main

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/kuzminykh.ulyana/task-3/internal/currency"
	"github.com/kuzminykh.ulyana/task-3/internal/decoder"
	"github.com/kuzminykh.ulyana/task-3/internal/models"
	"github.com/kuzminykh.ulyana/task-3/internal/saver"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unmarshaling config: %w", err)
	}

	currencies, err := decoder.DecodeFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("decoding file: %w", err)
	}

	sorted, err := currency.Sort(currencies.Valutes)
	if err != nil {
		return fmt.Errorf("sorting currencies: %w", err)
	}

	if err := saver.Save(sorted, cfg.OutputFile); err != nil {
		return fmt.Errorf("saving data: %w", err)
	}

	return nil
}
