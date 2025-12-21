package config

import (
	_ "embed"
)

var prodConfigData []byte

func loadProdConfig() ([]byte, error) {
	return prodConfigData, nil
}
