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
		case "utf-8", "utf8":
			return input, nil
		default:
			return nil, fmt.Errorf("unsupported charset %q", charset)
		}
	}

	if err := dec.Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	result := make([]Currency, 0, len(parsed.Valutes))
	for _, val := range parsed.Valutes {
		numCode := 0
		if s := strings.TrimSpace(val.NumCode); s != "" {
			if n, err := strconv.Atoi(s); err == nil {
				numCode = n
			}
		}

		valueStr := strings.ReplaceAll(strings.TrimSpace(val.Value), " ", "")
		valueStr = strings.ReplaceAll(valueStr, ",", ".")
		parsedValue, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse value %q: %w", valueStr, err)
		}

		nominal := 1.0
		if s := strings.TrimSpace(val.Nominal); s != "" {
			s2 := strings.ReplaceAll(s, ",", ".")
			if nf, err := strconv.ParseFloat(s2, 64); err == nil && nf > 0 {
				nominal = nf
			}
		}
		if nominal != 1 {
			parsedValue /= nominal
		}

		result = append(result, Currency{
			NumCode:  numCode,
			CharCode: strings.TrimSpace(val.CharCode),
			Value:    parsedValue,
		})
	}

	return result, nil
}
