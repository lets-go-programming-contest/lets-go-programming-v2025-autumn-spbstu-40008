package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidConfig = errors.New("invalid config: environment/log_level must be set")
	ErrUnmarshalYAML = errors.New("failed to unmarshal config YAML")
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(rawYAML, &cfg); err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrUnmarshalYAML, err)
	}

	if cfg.Environment == "" || cfg.LogLevel == "" {
		return Config{}, ErrInvalidConfig
	}

	return cfg, nil
}
