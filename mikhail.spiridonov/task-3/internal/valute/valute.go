package valute

import (
	"encoding/xml"
	"fmt"
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

	if err := dcdr.DecodeElement(&tempStrct, &start); err != nil {
		return fmt.Errorf("decode element: %w", err)
	}

	if tempStrct.NumCode == "" {
		strct.NumCode = 0
	} else {
		numCode, err := strconv.Atoi(tempStrct.NumCode)
		if err != nil {
			return fmt.Errorf("parse num code %q: %w", tempStrct.NumCode, err)
		}

		strct.NumCode = numCode
	}

	strct.CharCode = tempStrct.CharCode

	if tempStrct.Value == "" {
		strct.Value = 0

		return nil
	}

	normVal := strings.ReplaceAll(tempStrct.Value, ",", ".")
	value, err := strconv.ParseFloat(normVal, 64)
	if err != nil {
		return fmt.Errorf("parse value %q: %w", tempStrct.Value, err)
	}

	strct.Value = value

	return nil
}
