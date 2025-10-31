package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrMissingFields = errors.New("config must contain input-file and output-file")

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("reading YAML error: %w", err)
	}

	if config.InputFile == "" || config.OutputFile == "" {
		return nil, ErrMissingFields
	}

	return &config, nil
}
