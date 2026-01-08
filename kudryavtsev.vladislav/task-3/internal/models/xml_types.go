package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CustomFloat float64

func (cf *CustomFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var rawString string

	if err := d.DecodeElement(&rawString, &start); err != nil {
		return fmt.Errorf("element decode error: %w", err)
	}

	normalized := strings.ReplaceAll(rawString, ",", ".")

	parsedVal, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return fmt.Errorf("float conversion error: %w", err)
	}

	*cf = CustomFloat(parsedVal)

	return nil
}

type ExchangeData struct {
	Items []CurrencyItem `xml:"Valute"`
}

type CurrencyItem struct {
	NumericCode int         `json:"num_code"  xml:"NumCode"`
	LetterCode  string      `json:"char_code" xml:"CharCode"`
	Rate        CustomFloat `json:"value"     xml:"Value"`
}