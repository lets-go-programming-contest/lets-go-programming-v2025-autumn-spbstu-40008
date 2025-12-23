//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var prodConfig []byte

var cfg, _ = parseConfig(prodConfig)

func GetConfig() Config {
	return cfg
}
