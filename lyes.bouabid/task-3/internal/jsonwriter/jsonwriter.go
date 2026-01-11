package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveToJSON(filePath string, data interface{}, dirPerm os.FileMode, filePerm os.FileMode) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("serialize to JSON: %w", err)
	}

	directory := filepath.Dir(filePath)

	if err := os.MkdirAll(directory, dirPerm); err != nil {
		return fmt.Errorf("directory cannot being created '%s': %w", directory, err)
	}

	if err := os.WriteFile(filePath, jsonData, filePerm); err != nil {
		return fmt.Errorf("we cant write to file '%s': %w", filePath, err)
	}

	return nil
}
