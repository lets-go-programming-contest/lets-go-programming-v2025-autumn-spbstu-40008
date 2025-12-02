package currency

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type ExchangeRate float64

func (erv *ExchangeRate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return fmt.Errorf("ошибка чтения элемента: %w", err)
	}
	s = strings.Replace(s, ",", ".", 1)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("ошибка парсинга числа: %w", err)
	}
	*erv = ExchangeRate(f)
	return nil
}

type CurrencyIndex struct {
	XMLName    xml.Name `xml:"ValCurs"`
	Currencies []Item   `xml:"Valute"`
}

type Item struct {
	NumCode  int          `json:"num_code"  xml:"NumCode"`
	CharCode string       `json:"char_code" xml:"CharCode"`
	Value    ExchangeRate `json:"value"     xml:"Value"`
}
