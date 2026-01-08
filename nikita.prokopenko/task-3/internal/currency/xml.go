package currency

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type valCurs struct {
	Valutes []valute `xml:"Valute"`
}

type valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Value    string `xml:"Value"`
}

func DecodeXMLFile(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open xml %s: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	return DecodeXML(file)
}

func DecodeXML(r io.Reader) ([]Currency, error) {
	var parsed valCurs
	dec := xml.NewDecoder(r)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch strings.ToLower(strings.TrimSpace(charset)) {
		case "windows-1251", "cp1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return input, nil
		}
	}

	if err := dec.Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	result := make([]Currency, 0, len(parsed.Valutes))
	for _, val := range parsed.Valutes {
		numCodeStr := strings.TrimSpace(val.NumCode)
		if numCodeStr == "" {
			continue
		}
		numCode, err := strconv.Atoi(numCodeStr)
		if err != nil {
			continue
		}

		charCode := strings.TrimSpace(val.CharCode)

		nominalStr := strings.TrimSpace(val.Nominal)
		if nominalStr == "" {
			nominalStr = "1"
		}
		nominal, err := strconv.Atoi(nominalStr)
		if err != nil {
			return nil, fmt.Errorf("parse nominal: %w", err)
		}

		vStr := strings.ReplaceAll(strings.TrimSpace(val.Value), ",", ".")
		vFloat, err := strconv.ParseFloat(vStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse value: %w", err)
		}

		result = append(result, Currency{
			NumCode:  numCode,
			CharCode: charCode,
			Value:    vFloat / float64(nominal),
		})
	}

	return result, nil
}
