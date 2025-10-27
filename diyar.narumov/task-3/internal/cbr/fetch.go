package cbr

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/narumov-diyar/task-3/internal/currency"
	"golang.org/x/net/html/charset"
)

func Fetch(path string) ([]currency.Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			return
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var wrapper struct {
		Valutes []currency.Currency `xml:"Valute"`
	}

	if err := decoder.Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("unmarshal XML: %w", err)
	}

	currencies := wrapper.Valutes
	currency.SortByValueDesc(currencies)

	return currencies, nil
}
