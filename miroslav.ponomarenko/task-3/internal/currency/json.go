package currency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteJSON(list []Currency, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(outPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
