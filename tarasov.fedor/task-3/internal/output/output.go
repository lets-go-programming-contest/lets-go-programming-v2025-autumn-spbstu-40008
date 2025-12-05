package output

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateOutputFile(filename string) (*os.File, error) {
	dirPath := filepath.Dir(filename)

	const DirPerm = 0o755

	if err := os.MkdirAll(dirPath, DirPerm); err != nil {
		return nil, fmt.Errorf("unable to create output directory: %w", err)
	}

	const FilePerm = 0o644

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePerm)
	if err != nil {
		return nil, fmt.Errorf("unable to create output file: %w", err)
	}

	return file, nil
}
