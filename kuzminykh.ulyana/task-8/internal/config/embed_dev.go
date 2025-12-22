//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var devData []byte

func init() {
	initConfig(devData)
}
