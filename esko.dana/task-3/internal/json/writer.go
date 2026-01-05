package json

import (
	"encoding/json"
	"fmt"
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
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(outputPath, jsonData, filePermissions)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	return nil
}
