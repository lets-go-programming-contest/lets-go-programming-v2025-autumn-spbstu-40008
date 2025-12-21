package config

import (
	_ "embed"
)

var devConfigData []byte

func loadDevConfig() ([]byte, error) {
	return devConfigData, nil
}
