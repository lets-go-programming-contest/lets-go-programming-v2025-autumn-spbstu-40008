package parser

import (
	"errors"
	"fmt"
	"sort"

	"github.com/mordw1n/task-3/internal/jsonpack"
	"github.com/mordw1n/task-3/internal/valute"
	"github.com/mordw1n/task-3/internal/xmlpack"
)

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

	var ErrNoValidCurrencies = errors.New("no valid currencies with non-empty char code found")

	if len(validCurrencies) == 0 {
		return ErrNoValidCurrencies
	}

	sort.Slice(currencies, func(first, second int) bool {
		return currencies[first].Value > currencies[second].Value
	})

	if err := jsonpack.WriteInFile(outputFile, currencies); err != nil {
		return fmt.Errorf("write in JSON: %w", err)
	}

	return nil
}
