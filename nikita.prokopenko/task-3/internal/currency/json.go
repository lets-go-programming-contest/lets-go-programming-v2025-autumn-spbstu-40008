package currency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const dirPerm = 0o755

func SaveAsJSON(path string, items []Currency) error {
	dir := filepath.Dir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", path, err)
	}
	defer func() { _ = outFile.Close() }()

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "    ")
	if err := enc.Encode(items); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
