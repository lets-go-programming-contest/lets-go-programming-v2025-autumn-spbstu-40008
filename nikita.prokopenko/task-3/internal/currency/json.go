package currency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const directoryPerm = 0o750

func ExportToJSON(items []CurrencyItem, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, directoryPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("failed to close file: %v\n", closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(items); err != nil {
		return fmt.Errorf("failed to encode items: %w", err)
	}

	return nil
}
