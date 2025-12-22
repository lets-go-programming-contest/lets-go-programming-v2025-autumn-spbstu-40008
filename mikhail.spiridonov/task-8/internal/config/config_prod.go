//go:build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfig []byte

func loadConfig() Config {
	var cfg Config
	if err := yaml.Unmarshal(prodConfig, &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func init() {
	cfg = loadConfig()
}
