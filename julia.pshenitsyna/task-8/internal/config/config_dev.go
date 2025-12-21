package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

var devConfigData []byte

func init() {
	var cfg Config
	if err := yaml.Unmarshal(devConfigData, &cfg); err != nil {
		panic("failed to parse dev config: " + err.Error())
	}
	appConfig = &cfg
}