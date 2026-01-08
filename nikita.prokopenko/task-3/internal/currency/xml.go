package currency

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseCurrencyFile(filePath string) ([]CurrencyItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open currency file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("failed to close file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data CurrencyRates
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode xml: %w", err)
	}

	return data.Currencies, nil
}
