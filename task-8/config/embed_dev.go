//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var devData []byte

func Load() Config {
	return Get(devData)
}
