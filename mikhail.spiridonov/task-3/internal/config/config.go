package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Cfg struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadFile(filePath string) Cfg {
	data, _ := os.ReadFile(filePath)
	var cfg Cfg
	yaml.Unmarshal(data, &cfg)
	return cfg
}