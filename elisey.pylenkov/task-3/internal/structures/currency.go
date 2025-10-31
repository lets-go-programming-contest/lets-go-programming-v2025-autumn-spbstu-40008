package structures

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type CustomFloat float64

func (c *CustomFloat) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	s = strings.Replace(s, ",", ".", -1)

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*c = CustomFloat(f)
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
