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
	return LoadFromData(ConfigData)
}

func LoadFromData(data []byte) (Config, error) {
	var cfg Config
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("ошибка чтения конфига: %w", err)
	}
	return cfg, nil
}
