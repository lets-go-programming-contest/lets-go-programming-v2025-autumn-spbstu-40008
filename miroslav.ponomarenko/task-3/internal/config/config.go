package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(path string) (Config, error) {
	var cfg Config

	var ErrInvalidConfig = errors.New("invalid config fields")

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("open config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("decode yaml: %w", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return cfg, ErrInvalidConfig
	}

	return cfg, nil
}
