//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var prodConfig []byte

var cfg = func() Config { //nolint:gochecknoglobals
	c, err := parseConfig(prodConfig)
	if err != nil {
		panic(err)
	}

	return c
}()

func GetConfig() Config {
	return cfg
}
