package output

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func CreateFile(filename string) (*os.File, error) {
	dirPath := filepath.Dir(filename)

	if err := os.MkdirAll(dirPath, dirPerm); err != nil {
		return nil, fmt.Errorf("unable to create output directory: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return nil, fmt.Errorf("unable to create output file: %w", err)
	}

	return file, nil
}
