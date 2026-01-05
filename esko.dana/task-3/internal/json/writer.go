package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"esko.dana/task-3/internal/currency"
)

func Save(currencies []currency.Currency, outputPath string) error {
	outputDir := filepath.Dir(outputPath)
	err := os.MkdirAll(outputDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", outputDir, err)
	}

	jsonData, err := json.MarshalIndent(currencies, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	err = os.WriteFile(outputPath, jsonData, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write results to output file '%s': %w", outputPath, err)
	}

	return nil
}
