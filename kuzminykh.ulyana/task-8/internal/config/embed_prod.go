//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var configData []byte

var current, currentErr = New(configData)

func Get() Config {
	return current
}
