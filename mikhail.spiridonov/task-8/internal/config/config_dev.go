//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var devConfig []byte

var cfg, _ = parseConfig(devConfig)

func GetConfig() Config {
	return cfg
}
