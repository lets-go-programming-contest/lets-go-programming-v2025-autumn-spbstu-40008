//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var devConfig []byte

var cfg = func() Config {
	c, err := parseConfig(devConfig)
	if err != nil {
		panic(err)
	}
	return c
}()

func GetConfig() Config {
	return cfg
}
