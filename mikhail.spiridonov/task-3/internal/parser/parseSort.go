package parser

import (
	"sort"

	"github.com/mordw1n/task-3/internal/jsonpack"
	"github.com/mordw1n/task-3/internal/xmlpack"
)

func ParseAndSortXML(inputFile, outputFile string) {
	valCurs := xmlpack.ReadFile(inputFile)
	currencies := valCurs.Valutes
	
	sort.Slice(currencies, func(first, second int) bool {
		return currencies[first].Value > currencies[second].Value
	})
	
	for index, currency := range currencies {
		if currency.NumCode == 0 {
			return fmt.Errorf("Invalid num code")
		}
		if currency.CharCode == "" {
			return fmt.Errorf("Empty char code")
		}
		if currency.Value <= 0 {
			return fmt.Errorf("Bad value of valute")
		}
		if index > 0 && currencies[index-1].Value < currency.Value {
			return fmt.Errorf("Incorrect sort")
		}
	}

	if err := jsonpack.WriteInFile(outputFile, currencies); err != nil {
		return fmt.Errorf("Write in JSON: %w", err)
	}

	return nil
}