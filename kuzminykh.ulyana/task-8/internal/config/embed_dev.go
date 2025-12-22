//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var configData []byte

var current, currentErr = New(configData)

func Get() Config {
	return current
}
