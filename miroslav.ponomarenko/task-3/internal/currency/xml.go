package currency

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ReadXML(path string) (*ExchangeRates, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open xml: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	dec := xml.NewDecoder(file)
	dec.CharsetReader = charset.NewReaderLabel

	var rates ExchangeRates

	if err := dec.Decode(&rates); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &rates, nil
}
