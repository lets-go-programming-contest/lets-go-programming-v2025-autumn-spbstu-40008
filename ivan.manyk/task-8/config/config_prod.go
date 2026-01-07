// config/config_prod.go
//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var prodData []byte

func LoadConf() Conf {
	return Get(prodData)
}