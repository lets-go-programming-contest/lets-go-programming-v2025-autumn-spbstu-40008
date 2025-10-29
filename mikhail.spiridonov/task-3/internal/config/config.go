package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadFile(filePath string) Config {
	data, _ := os.ReadFile(filePath)
	var config Config
	yaml.Unmarshal(data, &config)
	return config
}