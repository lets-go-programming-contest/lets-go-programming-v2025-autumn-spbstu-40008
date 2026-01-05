package saver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuzminykh.ulyana/task-3/internal/models"
)

func Save(data []models.Output, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	return nil
}
