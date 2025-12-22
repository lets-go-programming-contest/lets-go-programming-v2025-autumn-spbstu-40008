//go:build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfig []byte

func init() {
	if err := yaml.Unmarshal(prodConfig, &cfg); err != nil {
		panic(err)
	}
}
