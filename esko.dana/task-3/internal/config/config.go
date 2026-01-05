package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputPath    string `yaml:"inputPath"`
	CurrencyCode string `yaml:"currencyCode"`
	OutputPath   string `yaml:"outputPath"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
