package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"

	"task-3/internal/structures"
)

var (
	ErrUnsupportedEncoding = errors.New("unsupported encoding")
	ErrNoValuteData        = errors.New("xml doesn't contain valute data")
)

func ParseCurrencyXML(filePath string) (*structures.ValCurs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening xml file %s: %w", filePath, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedEncoding, charset)
		}
	}

	var valCurs structures.ValCurs
	err = decoder.Decode(&valCurs)
	if err != nil {
		return nil, fmt.Errorf("error reading xml: %w", err)
	}

	if len(valCurs.Valutes) == 0 {
		return nil, ErrNoValuteData
	}

	return &valCurs, nil
}

func ConvertToOutput(valutes []structures.Valute) []structures.OutputCurrency {
	output := make([]structures.OutputCurrency, 0, len(valutes))

	for _, valute := range valutes {
		output = append(output, structures.OutputCurrency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    float64(valute.Value),
		})
	}

	return output
}
