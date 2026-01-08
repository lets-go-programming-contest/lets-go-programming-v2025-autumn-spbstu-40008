package appconfig

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	SourceFile string `yaml:"input-file"`
	TargetFile string `yaml:"output-file"`
}

func New(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Settings
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
