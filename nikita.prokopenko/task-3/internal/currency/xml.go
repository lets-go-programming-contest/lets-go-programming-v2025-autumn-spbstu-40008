package currency

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseCurrencyFile(filePath string) ([]CurrencyItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data CurrencyRates
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data.Currencies, nil
}
