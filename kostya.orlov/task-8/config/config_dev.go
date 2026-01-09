//go:build dev
// +build dev

package config

import _ "embed"

//go:embed dev.yaml
var devYAML []byte

func rawYAML() []byte { return devYAML }
