//go:build !dev

package config

import (
	_ "embed"
)

//go:embed prod.yaml
var ConfigData []byte

var Current, CurrentErr = New(ConfigData)
