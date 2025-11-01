package currency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	dirPerm  os.FileMode = 0o755
	filePerm os.FileMode = 0o600
)

func WriteJSON(list []Currency, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), dirPerm); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	data = append(data, '\n')

	if err := os.WriteFile(outPath, data, filePerm); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
