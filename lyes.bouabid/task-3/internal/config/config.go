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

func LoadConfig(configPath string) (*Config, error) {
	var result Config

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("we cant read the config file: %w", err)
	}

	err = yaml.Unmarshal(file, &result)
	if err != nil {
		return nil, fmt.Errorf("we cant unmarshal file: %w", err)
	}

	return &result, nil
}
