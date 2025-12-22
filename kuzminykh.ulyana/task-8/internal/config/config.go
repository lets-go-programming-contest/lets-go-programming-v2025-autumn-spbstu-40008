package config

import (
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

var cfg Config

func Get() Config {
	return cfg
}

func initConfig(data []byte) {
	yaml.Unmarshal(data, &cfg)
}
