package parser

import (
	"errors"
	"fmt"
	"sort"

	"github.com/mordw1n/task-3/internal/jsonpack"
	"github.com/mordw1n/task-3/internal/xmlpack"
)

var errNoValidCurrencies = errors.New("no valid currencies found")

func ParseAndSortXML(inputFile, outputFile string) error {
	valCurs, err := xmlpack.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("read XML file: %w", err)
	}

	currencies := valCurs.Valutes

	if len(currencies) == 0 {
		return errNoValidCurrencies
	}

	sort.Slice(currencies, func(first, second int) bool {
		return currencies[first].Value > currencies[second].Value
	})

	if err := jsonpack.WriteInFile(outputFile, currencies); err != nil {
		return fmt.Errorf("write in JSON: %w", err)
	}

	return nil
}
