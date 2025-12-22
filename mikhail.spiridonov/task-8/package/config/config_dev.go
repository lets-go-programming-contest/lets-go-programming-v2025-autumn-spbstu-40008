//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var devConfig []byte

func init() {
	if err := yaml.Unmarshal(devConfig, &cfg); err != nil {
		panic(err)
	}
}
