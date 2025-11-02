package app

import (
	"fmt"

	"github.com/narumov-diyar/task-3/internal/cbr"
	"github.com/narumov-diyar/task-3/internal/config"
	"github.com/narumov-diyar/task-3/internal/format"
)

func Run(configPath, outputFormat string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	currencies, err := cbr.Fetch(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("parse XML: %w", err)
	}

	if err := format.Write(currencies, cfg.OutputFile, outputFormat); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
