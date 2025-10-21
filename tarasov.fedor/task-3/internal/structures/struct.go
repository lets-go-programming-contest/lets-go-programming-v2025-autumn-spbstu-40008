package structures

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type CustomFloat float64

func (cfg *CustomFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.ReplaceAll(s, ",", ".")

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*cfg = CustomFloat(val)
	return nil
}

type File struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

type ValCurs struct {
	Valute []Valute `xml:"Valute" json:"-"`
}
type Valute struct {
	NumCode  int         `xml:"NumCode" json:"num_code"`
	CharCode string      `xml:"CharCode" json:"char_code"`
	Value    CustomFloat `xml:"Value" json:"value"`
}
