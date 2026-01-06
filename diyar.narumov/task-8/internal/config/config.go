package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	ErrUnmarshal = errors.New("failed to unmarshal config YAML")
	ErrInvalid   = errors.New("invalid config: environment/log_level must be set")
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(rawYAML, &cfg); err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrUnmarshal, err)
	}

	if cfg.Environment == "" || cfg.LogLevel == "" {
		return Config{}, ErrInvalid
	}

	return cfg, nil
}
