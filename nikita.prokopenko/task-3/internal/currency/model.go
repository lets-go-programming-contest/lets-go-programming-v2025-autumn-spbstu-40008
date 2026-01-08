package conversion

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

type ExchangeData struct {
	Items []CurrencyItem `xml:"Valute"`
}

func (e *ExchangeData) GetAllCurrencies() []CurrencyItem {
	return e.Items
}

type CurrencyItem struct {
	NumericCode int     `json:"num_code" xml:"NumCode"`
	AlphaCode   string  `json:"char_code" xml:"CharCode"`
	Value       Decimal `json:"value" xml:"Value"`
}