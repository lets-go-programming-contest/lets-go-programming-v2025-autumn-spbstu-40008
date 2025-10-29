package parser

import (
	"errors"
	"fmt"
	"sort"

	"github.com/mordw1n/task-3/internal/jsonpack"
	"github.com/mordw1n/task-3/internal/valute"
	"github.com/mordw1n/task-3/internal/xmlpack"
)

var errNoValidCurrencies = errors.New("no valid currencies with non-empty char code found")

func ParseAndSortXML(inputFile, outputFile string) error {
	valCurs, err := xmlpack.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("read XML file: %w", err)
	}

	currencies := valCurs.Valutes

	validCurrencies := make([]valute.StructOfXMLandJSON, 0, len(currencies))

	for _, currency := range currencies {
		if currency.CharCode != "" {
			validCurrencies = append(validCurrencies, currency)
		}
	}

	if len(validCurrencies) == 0 {
		return errNoValidCurrencies
	}

	sort.Slice(validCurrencies, func(first, second int) bool {
		return validCurrencies[first].Value > validCurrencies[second].Value
	})

	if err := jsonpack.WriteInFile(outputFile, validCurrencies); err != nil {
		return fmt.Errorf("write in JSON: %w", err)
	}

	return nil
}
