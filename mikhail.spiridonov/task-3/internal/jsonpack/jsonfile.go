package jsonpack

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mordw1n/task-3/internal/valute"
)

func WriteInFile(filePath string, currencies []valute.StructOfXMLandJSON) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("Create directory for %q: %w", filePath, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Create JSON %q: %w", filePath, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(fmt.Errorf("Close JSON %q: %w", filePath, closeErr))
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(currencies); err != nil {
		return fmt.Errorf("Encode to JSON %q: %w", filePath, err)
	}

	return nil
}
