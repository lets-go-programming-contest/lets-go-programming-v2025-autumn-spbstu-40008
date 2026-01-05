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
	if filePath == "" {
		return nil, fmt.Errorf("failed to open XML file '': open : no such file or directory")
	}

	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file '%s': %w", filePath, err)
	}

	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs

	err = decoder.Decode(&valCurs)

	closeErr := xmlFile.Close()
	if closeErr != nil {
		return nil, fmt.Errorf("failed to close XML file: %w", closeErr)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode XML from '%s': %w", filePath, err)
	}

	parsedValutes := make([]ParsedValute, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		valueStr := strings.ReplaceAll(valute.Value, ",", ".")

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
