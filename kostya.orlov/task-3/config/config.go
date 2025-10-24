package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)

	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var config Config

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &config, nil
}
