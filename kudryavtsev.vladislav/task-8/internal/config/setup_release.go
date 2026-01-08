//go:build !dev

package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed settings-prod.yaml

var rawProdSource []byte

func Load() (*Settings, error) {

	var appParams Settings

	if parseErr := yaml.Unmarshal(rawProdSource, &appParams); parseErr != nil {

		return nil, fmt.Errorf("критическая ошибка загрузки конфига: %w", parseErr)

	}

	return &appParams, nil

}
