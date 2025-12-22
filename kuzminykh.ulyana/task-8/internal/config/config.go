package config

import "gopkg.in/yaml.v3"

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func New(data []byte) (Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
