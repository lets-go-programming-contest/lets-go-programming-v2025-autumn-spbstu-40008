package parser

import (
	"sort"
	
	"github.com/mordw1n/task-3/internal/jsonpack"
	"github.com/mordw1n/task-3/internal/xmlpack"
)

func ParseAndSortXML(inputFile, outputFile string) {
	valCurs := xmlpack.ReadFile(inputFile)
	currencies := valCurs.Valutes
	
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
	
	jsonpack.WriteInFile(outputFile, currencies)
}