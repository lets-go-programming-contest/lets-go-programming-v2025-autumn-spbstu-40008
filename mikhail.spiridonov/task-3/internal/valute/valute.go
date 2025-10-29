package valute

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type StructOfXMLandJSON struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

type ValuteCurs struct {
	XMLName xml.Name             `xml:"ValCurs"`
	Valutes []StructOfXMLandJSON `xml:"Valute"`
}

func (strct *StructOfXMLandJSON) UnmarshalXML(dcdr *xml.Decoder, start xml.StartElement) error {
	type temp struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}
	var tempStrct temp
	dcdr.DecodeElement(&tempStrct, &start)

	strct.NumCode, _ = strconv.Atoi(tempStrct.NumCode)
	strct.CharCode = tempStrct.CharCode
	normVal := strings.ReplaceAll(tempStrct.Value, ",", ".")
	strct.Value, _ = strconv.ParseFloat(normVal, 64)

	return nil
}
