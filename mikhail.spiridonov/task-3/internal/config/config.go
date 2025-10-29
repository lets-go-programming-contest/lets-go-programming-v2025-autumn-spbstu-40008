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

func LoadFile(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("Read XML %q: %w", filePath, err)
	}
	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("Unmarshal XML %q: %w", filePath, err)
	}

	return config, nil
}