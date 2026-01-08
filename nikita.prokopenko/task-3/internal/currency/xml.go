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
	defer file.Close()
	return DecodeXML(file)
}

func DecodeXML(r io.Reader) ([]Currency, error) {
	dec := xml.NewDecoder(r)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}

	var parsed valCurs
	if err := dec.Decode(&parsed); err != nil {
		return nil, err
	}

	var result []Currency
	for _, v := range parsed.Valutes {
		num, _ := strconv.Atoi(v.NumCode)
		nom, _ := strconv.ParseFloat(strings.ReplaceAll(v.Nominal, ",", "."), 64)
		val, _ := strconv.ParseFloat(strings.ReplaceAll(v.Value, ",", "."), 64)

		if nom == 0 {
			nom = 1
		}

		result = append(result, Currency{
			NumCode:  num,
			CharCode: v.CharCode,
			Value:    val / nom,
		})
	}
	return result, nil
}
