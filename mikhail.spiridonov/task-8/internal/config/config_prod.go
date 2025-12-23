//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var prodConfig []byte

func getEmbeddedConfig() []byte {
	return prodConfig
}
