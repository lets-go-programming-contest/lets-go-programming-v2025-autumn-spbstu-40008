package config

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Get(data []byte) Config {
	var cfg Config
	_ = yaml.Unmarshal(data, &cfg)
	return cfg
}
