package config

import (
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

var (
	config     Config
	configOnce sync.Once
	configErr  error
)

func GetConfig() (Config, error) {
	configOnce.Do(func() {
		config, configErr = loadConfig()
	})

	if configErr != nil {
		return Config{}, configErr
	}
    
	return config, nil
}

func GetConfigOrPanic() Config {
	cfg, err := GetConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadConfig() (Config, error) {
	var data []byte

	data = getEmbeddedConfig()

	return parseConfig(data)
}

func parseConfig(data []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal YAML config: %w", err)
	}

	return cfg, nil
}
