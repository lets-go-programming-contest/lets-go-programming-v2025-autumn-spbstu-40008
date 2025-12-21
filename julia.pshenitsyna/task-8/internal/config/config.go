package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (Config, error) {
	return parseConfig(ConfigData)
}

func parseConfig(data []byte) (Config, error) {
	var conf Config
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, fmt.Errorf("Config reading error: %w", err)
	}

	return conf, nil
}
