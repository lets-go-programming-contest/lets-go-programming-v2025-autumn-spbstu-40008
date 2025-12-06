package structures

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CustomFloat float64

func (c *CustomFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var stringValue string

	if err := d.DecodeElement(&stringValue, &start); err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	stringValue = strings.ReplaceAll(stringValue, ",", ".")

	floatValue, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return fmt.Errorf("parse float: %w", err)
	}

	*c = CustomFloat(floatValue)

	return nil
}

type Valute struct {
	NumCode  int         `xml:"NumCode"`
	CharCode string      `xml:"CharCode"`
	Value    CustomFloat `xml:"Value"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type OutputCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}
