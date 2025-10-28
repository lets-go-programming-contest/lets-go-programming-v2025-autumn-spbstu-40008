package valute

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type StructOfXMLandJSON struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

type ValuteCurs struct {
	XMLName xml.Name               `xml:"ValCurs"`
	Valutes []StructOfXMLandJSON   `xml:"Valute"`
}

func (s *StructOfXMLandJSON) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type temp struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}
	var t temp
	d.DecodeElement(&t, &start)
	
	s.NumCode, _ = strconv.Atoi(t.NumCode)
	s.CharCode = t.CharCode
	s.Value, _ = strconv.ParseFloat(strings.Replace(t.Value, ",", ".", -1), 64)
	return nil
}