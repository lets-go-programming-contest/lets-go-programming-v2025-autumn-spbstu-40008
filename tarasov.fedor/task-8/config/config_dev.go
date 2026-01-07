//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var yamlFile []byte

func GetConfig() Config {
	return ParseConfig(yamlFile)
}
