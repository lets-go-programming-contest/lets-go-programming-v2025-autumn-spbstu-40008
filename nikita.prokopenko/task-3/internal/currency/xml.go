package conversion

import (
	"encoding/xml"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseCurrencyData(filePath string) (*ExchangeData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data ExchangeData
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}