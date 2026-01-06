package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidConfig = errors.New("invalid configuration")
	ErrEmptyFields   = errors.New("environment and log_level must be specified")
)

type AppConfig struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Parse(data []byte) (*AppConfig, error) {
	var cfg AppConfig

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}

	if cfg.Environment == "" || cfg.LogLevel == "" {
		return nil, ErrEmptyFields
	}

	return &cfg, nil
}
