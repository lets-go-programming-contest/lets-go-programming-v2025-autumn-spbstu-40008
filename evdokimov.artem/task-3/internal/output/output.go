package output

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(filename string) (*os.File, error) {
	dirPath := filepath.Dir(filename)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, fmt.Errorf("unable to create output directory: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to create output file: %w", err)
	}

	return file, nil
}
