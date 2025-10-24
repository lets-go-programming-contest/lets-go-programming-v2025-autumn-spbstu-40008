package currency

import "encoding/xml"

type ValuteXML struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type ValCurs struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Valutes []ValuteXML `xml:"Valute"`
}

type ResultValute struct {
	NumCode  int     `json:"num_code"  yaml:"num_code"  xml:"num_code"`
	CharCode string  `json:"char_code" yaml:"char_code" xml:"char_code"`
	Value    float64 `json:"value"     yaml:"value"     xml:"value"`
}
