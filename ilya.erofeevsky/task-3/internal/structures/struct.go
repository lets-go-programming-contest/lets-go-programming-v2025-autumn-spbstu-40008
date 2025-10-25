package structures

import (
	"encoding/xml"
)

type ValuteXML struct {
	NumCode  string `xml: "NumCode"`
	CharCode string `xml: "CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type ReadingXML struct {
	XMLName     xml.Name    `xml:"ValCurs"`
	Information []ValuteXML `xml:"Valute"`
}

type File struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

type ProcessedCurrency struct {
	NumCode  int
	CharCode string
	Value    float64
	Nominal  int
}

type ResultItem struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}
