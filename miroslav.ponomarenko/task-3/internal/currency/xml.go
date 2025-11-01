package currency

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ReadXML(path string) (*ExchangeRates, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open xml: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	dec := xml.NewDecoder(f)
	dec.CharsetReader = charset.NewReaderLabel

	var rates ExchangeRates
	if err := dec.Decode(&rates); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &rates, nil
}
