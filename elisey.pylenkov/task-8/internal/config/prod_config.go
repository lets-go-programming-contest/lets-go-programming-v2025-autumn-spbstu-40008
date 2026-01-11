//go:build !dev

package config

import _ "embed"

//go:embed prod.yaml
var prodConfigContent []byte

func Load() (*AppConfig, error) {
	return Parse(prodConfigContent)
}
