//go:build dev

package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed settings-dev.yaml
var rawDevSource []byte

func Load() (*Settings, error) {
	var appParams Settings

	if parseErr := yaml.Unmarshal(rawDevSource, &appParams); parseErr != nil {
		return nil, fmt.Errorf("сбой при разборе dev-конфига: %w", parseErr)
	}

	return &appParams, nil
}
