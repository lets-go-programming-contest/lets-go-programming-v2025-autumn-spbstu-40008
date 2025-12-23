package config

import (
	"gopkg.in/yaml.v3"
)

type Conf struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Get(data []byte) Conf {
	var cfg Conf
	_ = yaml.Unmarshal(data, &cfg)

	return cfg
}
