package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gagysun/task-3/internal/models"
	"golang.org/x/net/html/charset"
)

const filePerm = 0o600

func LoadXML(filePath string) (*models.ExchangeData, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read source file %s: %w", filePath, err)
	}

	var result models.ExchangeData

	reader := bytes.NewReader(fileContent)

	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("xml parsing error: %w", err)
	}

	return &result, nil
}

func ExportToJSON(data []models.CurrencyItem, path string) error {
	encodedJSON, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("serialization error: %w", err)
	}

	directory := filepath.Dir(path)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("filesystem error (mkdir): %w", err)
	}

	if err := os.WriteFile(path, encodedJSON, filePerm); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}
