//go:build dev

package config

import _ "embed"

//go:embed dev.yaml
var devConfigContent []byte

func Load() (*AppConfig, error) {
	return Parse(devConfigContent)
}
