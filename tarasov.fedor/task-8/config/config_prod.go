//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var yamlFile []byte

func GetConfig() Config {
	return ParseConfig(yamlFile)
}
