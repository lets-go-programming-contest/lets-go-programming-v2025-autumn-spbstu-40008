package config

import (
	"gopkg.in/yaml.v3"
)

type Cfg struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadFile() {
	
}