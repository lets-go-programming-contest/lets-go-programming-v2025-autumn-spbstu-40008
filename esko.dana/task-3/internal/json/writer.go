package json

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	dirPermissions  = 0o755
	filePermissions = 0o600
)

func Write(data interface{}, outputPath string) error {
	outputDir := filepath.Dir(outputPath)
	err := os.MkdirAll(outputDir, dirPermissions)

	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath, jsonData, filePermissions)

	return err
}
