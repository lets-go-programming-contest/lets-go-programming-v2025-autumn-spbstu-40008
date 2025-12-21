package config

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var configData []byte
	var err error

	configData, err = loadDevConfig()
	if err == nil {
		return parseConfig(configData)
	}

	configData, err = loadProdConfig()
	if err != nil {
		return nil, err
	}

	return parseConfig(configData)
}

func parseConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
