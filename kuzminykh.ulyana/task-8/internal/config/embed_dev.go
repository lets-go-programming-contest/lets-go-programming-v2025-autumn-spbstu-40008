//go:build dev

package config

import (
	_ "embed"
)

//go:embed dev.yaml
var ConfigData []byte

var Current, CurrentErr = New(ConfigData)
