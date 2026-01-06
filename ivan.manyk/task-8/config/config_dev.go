// config/config_dev.go
//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var devData []byte

func LoadConf() Conf {
	return Get(devData)
}