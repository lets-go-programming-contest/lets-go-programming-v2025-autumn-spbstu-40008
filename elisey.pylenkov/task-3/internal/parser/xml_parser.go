package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"task-3/internal/structures"

	"golang.org/x/text/encoding/charmap"
)

func ParseCurrencyXML(filePath string) (*structures.ValCurs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening xml file %s: %w", filePath, err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unsupported encoding %s", charset)
		}
	}

	var valCurs structures.ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("error reading xml: %w", err)
	}

	if len(valCurs.Valutes) == 0 {
		return nil, fmt.Errorf("xml doesn't contain valute data")
	}

	return &valCurs, nil
}

func ConvertToOutput(valutes []structures.Valute) []structures.OutputCurrency {
	var output []structures.OutputCurrency

	for _, valute := range valutes {
		output = append(output, structures.OutputCurrency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    float64(valute.Value),
		})
	}

	return output
}
