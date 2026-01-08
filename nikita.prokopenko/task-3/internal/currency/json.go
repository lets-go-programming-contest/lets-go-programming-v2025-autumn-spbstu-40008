package currency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveAsJSON(path string, items []Currency) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() { _ = outFile.Close() }()

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "    ")
	if err := enc.Encode(items); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
