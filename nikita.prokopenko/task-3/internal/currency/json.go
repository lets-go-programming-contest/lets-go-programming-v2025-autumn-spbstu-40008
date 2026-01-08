package currency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveAsJSON(path string, items []Currency) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	return enc.Encode(items)
}
