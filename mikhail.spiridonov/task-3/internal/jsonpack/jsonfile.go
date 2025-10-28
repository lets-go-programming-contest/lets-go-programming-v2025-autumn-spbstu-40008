package jsonpack

import (
	"os"
	"path/filepath"
	"encoding/json"
	"github.com/mordw1n/task-3/valute"
)

func WriteInFile(filePath string, currencies []valute.StructOfXMLandJSON) {
	os.MkdirAll(filepath.Dir(filePath), 0755)
	file, _ := os.Create(filePath)
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encoder.Encode(currencies)
}