package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

var ErrEmptyRequiredFields = errors.New("config: empty required fields")

func Load() (Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(rawYAML(), &cfg); err != nil {
		return Config{}, fmt.Errorf("config: unmarshal yaml: %w", err)
	}

	if cfg.Environment == "" || cfg.LogLevel == "" {
		return Config{}, ErrEmptyRequiredFields
	}

	return cfg, nil
}
