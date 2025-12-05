package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type ParsedValute struct {
	NumCode  string
	CharCode string
	Value    float64
}

func Parse(filePath string) ([]ParsedValute, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file '%s': %w", filePath, err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode XML from '%s': %w", filePath, err)
	}

	var parsedValutes []ParsedValute
	for _, valute := range valCurs.Valutes {
		valueStr := strings.Replace(valute.Value, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid 'Value' format '%s' in XML: %w", valute.Value, err)
		}

		parsedValutes = append(parsedValutes, ParsedValute{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	return parsedValutes, nil
}
