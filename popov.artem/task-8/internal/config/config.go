package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(yamlData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config error: %w", err)
	}
	return &cfg, nil
}
