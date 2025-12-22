//go:build dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfig []byte

func loadConfig() Config {
	var cfg Config
	if err := yaml.Unmarshal(devConfig, &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func init() {
	cfg = loadConfig()
}
