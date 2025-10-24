package structures

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
)

type CustomFloat float64

var (
	ErrDecodeXML  = errors.New("failed to decode XML element into string")
	ErrParseFloat = errors.New("failed to parse float from normalized string")
)

func (cfg *CustomFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valueStr string

	if err := d.DecodeElement(&valueStr, &start); err != nil {
		return ErrDecodeXML
	}

	valueStr = strings.TrimSpace(valueStr)
	valueStr = strings.ReplaceAll(valueStr, ",", ".")

	val, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return ErrParseFloat
	}

	*cfg = CustomFloat(val)

	return nil
}

type ValCurs struct {
	Valute []Valute `json:"valute" xml:"Valute"`
}
type Valute struct {
	NumCode  int         `json:"num_code"  xml:"NumCode"`
	CharCode string      `json:"char_code" xml:"CharCode"`
	Value    CustomFloat `json:"value"     xml:"Value"`
}
