package utils

import (
	"sort"

	"student/currency-processor/internal/models"
)

func SortCurrencyData(source *models.ExchangeData) []models.CurrencyItem {
	output := make([]models.CurrencyItem, len(source.Items))

	copy(output, source.Items)

	sort.Slice(output, func(i, j int) bool {
		return output[i].Rate > output[j].Rate
	})

	return output
}