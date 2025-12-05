package structures

import (
	"encoding/xml"
)

type Currency struct {
	ID         string `json:"-"         xml:"ID,attr"`
	NumCodeStr string `json:"-"         xml:"NumCode"`
	CharCode   string `json:"char_code" xml:"CharCode"`
	NominalStr string `json:"-"         xml:"Nominal"`
	Name       string `json:"-"         xml:"Name"`
	ValueStr   string `json:"-"         xml:"Value"`

	NumCode int     `json:"num_code" xml:"-"`
	Value   float64 `json:"value"    xml:"-"`
}

type ReadingXML struct {
	XMLName     xml.Name   `xml:"ValCurs"`
	Information []Currency `xml:"Valute"`
}
