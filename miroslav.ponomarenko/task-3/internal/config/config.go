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

func Load(path string) (c Config, err error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("Open config: %w", err)
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("Decode yaml: %w", err)
	}

	if c.InputFile == "" || c.OutputFile == "" {
		return c, fmt.Errorf("Invalid config fields")
	}

	return c, nil
}
