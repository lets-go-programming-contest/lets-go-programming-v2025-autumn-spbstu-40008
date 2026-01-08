package currency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ExportToJSON(items []CurrencyItem, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
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
	
	return encoder.Encode(items)
}
