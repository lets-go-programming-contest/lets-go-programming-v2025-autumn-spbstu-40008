package structures

import (
	"encoding/xml"
)

type Currency struct {
	ID         string `xml:"ID,attr" json:"-"`
	NumCodeStr string `xml:"NumCode" json:"-"`
	CharCode   string `xml:"CharCode" json:"char_code"`
	NominalStr string `xml:"Nominal" json:"-"`
	Name       string `xml:"Name" json:"-"`
	ValueStr   string `xml:"Value" json:"-"`

	NumCode int     `xml:"-" json:"num_code"`
	Value   float64 `xml:"-" json:"value"`
}

type ReadingXML struct {
	XMLName     xml.Name   `xml:"ValCurs"`
	Information []Currency `xml:"Valute"`
}

type File struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}
