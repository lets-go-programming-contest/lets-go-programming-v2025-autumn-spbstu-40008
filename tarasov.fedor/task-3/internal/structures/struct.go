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

type File struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

type ValCurs struct {
	Valute []Valute `xml:"Valute"  json:"valute"`
}
type Valute struct {
	NumCode  int         `xml:"NumCode"  json:"num_code"`
	CharCode string      `xml:"CharCode" json:"char_code"`
	Value    CustomFloat `xml:"Value"    json:"value"`
}
