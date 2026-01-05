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
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	encodeErr := encoder.Encode(data)
	closeErr := file.Close()

	if encodeErr != nil {
		return fmt.Errorf("encoding JSON: %w", encodeErr)
	}

	if closeErr != nil {
		return fmt.Errorf("closing file: %w", closeErr)
	}

	return nil
}
