package config

import (
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
		return File{}, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return File{}, err
	}

	return cfg, nil
}
