package currency

import (
	"strconv"
	"strings"
)

type Decimal float64

func (d *Decimal) UnmarshalText(text []byte) error {
	cleanText := strings.ReplaceAll(strings.TrimSpace(string(text)), ",", ".")
	value, err := strconv.ParseFloat(cleanText, 64)
	if err != nil {
		return err
	}
	
	*d = Decimal(value)

	return nil
}

type CurrencyRates struct {
	Currencies []CurrencyItem `xml:"Valute"`
}

type CurrencyItem struct {
	NumCode  int     `json:"num_code" xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    Decimal `json:"value" xml:"Value"`
}
