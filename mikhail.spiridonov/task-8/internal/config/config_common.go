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

type configManager struct {
	config     Config
	configOnce sync.Once
	configErr  error
}

var manager = &configManager{}

func GetConfig() (Config, error) {
	manager.configOnce.Do(func() {
		data := getEmbeddedConfig()
		manager.config, manager.configErr = parseConfig(data)
	})

	if manager.configErr != nil {
		return Config{}, manager.configErr
	}

	return manager.config, nil
}

func GetConfigOrPanic() Config {
	cfg, err := GetConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

func parseConfig(data []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal YAML config: %w", err)
	}

	return cfg, nil
}

func getEmbeddedConfig() []byte {
	return nil
}
