package currency

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Currency struct {
	NumCode  int     `json:"num_code"  yaml:"num_code"`
	CharCode string  `json:"char_code" yaml:"char_code"`
	Value    float64 `json:"value"     yaml:"value"`
}

func (c *Currency) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var temp struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	if err := dec.DecodeElement(&temp, &start); err != nil {
		return fmt.Errorf("decode XML element: %w", err)
	}

	if temp.NumCode == "" {
		temp.NumCode = "0"
	}

	num, err := strconv.Atoi(temp.NumCode)
	if err != nil {
		return fmt.Errorf("invalid NumCode %q: %w", temp.NumCode, err)
	}

	valueStr := strings.ReplaceAll(temp.Value, ",", ".")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return fmt.Errorf("invalid Value %q: %w", temp.Value, err)
	}

	c.NumCode = num
	c.CharCode = temp.CharCode
	c.Value = value

	return nil
}

func (c Currency) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	numStr := fmt.Sprintf("%03d", c.NumCode)
	valueStr := strings.ReplaceAll(fmt.Sprintf("%.4f", c.Value), ".", ",")

	if err := enc.EncodeElement(struct {
		XMLName  xml.Name `xml:"Valute"`
		NumCode  string   `xml:"NumCode"`
		CharCode string   `xml:"CharCode"`
		Value    string   `xml:"Value"`
	}{
		XMLName:  xml.Name{Space: "", Local: "Valute"},
		NumCode:  numStr,
		CharCode: c.CharCode,
		Value:    valueStr,
	}, start); err != nil {
		return fmt.Errorf("encode XML element: %w", err)
	}

	return nil
}

func SortByValueDesc(currencies []Currency) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
}
