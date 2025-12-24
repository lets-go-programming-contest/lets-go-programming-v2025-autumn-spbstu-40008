//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var prodData []byte

func Load() Conf {
	return Get(prodData)
}
