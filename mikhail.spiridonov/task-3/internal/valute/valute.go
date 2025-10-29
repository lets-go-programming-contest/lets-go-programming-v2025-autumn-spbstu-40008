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
	strct.Value, _ = strconv.ParseFloat(strings.Replace(tempStrct.Value, ",", ".", -1), 64)
	return nil
}