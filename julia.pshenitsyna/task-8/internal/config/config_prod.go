package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

var prodConfigData []byte

func init() {
	var cfg Config
	if err := yaml.Unmarshal(prodConfigData, &cfg); err != nil {
		panic("failed to parse prod config: " + err.Error())
	}
	appConfig = &cfg
}