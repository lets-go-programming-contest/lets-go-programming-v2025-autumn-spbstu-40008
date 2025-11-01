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

func Load(path string) (Config, error) {
	var c Config

	data, err := os.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("open config: %w", err)
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("decode yaml: %w", err)
	}

	if c.InputFile == "" || c.OutputFile == "" {
		return c, fmt.Errorf("invalid config fields")
	}

	return c, nil
}
