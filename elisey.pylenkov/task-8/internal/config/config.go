package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Parse(data []byte) (*AppConfig, error) {
	var cfg AppConfig

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("configuration parsing failed: %w", err)
	}

	if cfg.Environment == "" || cfg.LogLevel == "" {
		return nil, fmt.Errorf("environment and log_level must be specified")
	}

	return &cfg, nil
}
