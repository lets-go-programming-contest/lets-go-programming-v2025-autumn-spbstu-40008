//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var prodConfig []byte

func Load() Config {
	return parse(prodConfig)
}
