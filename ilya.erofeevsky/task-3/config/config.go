package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type File struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

func ReadFile(configPath string) (File, error) {
	var cfg File

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return File{}, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return File{}, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return cfg, nil
}
