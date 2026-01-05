package currency

import (
	"fmt"
	"sort"
	"strconv"

	"esko.dana/task-3/internal/xml"
)

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type ByValueDesc []Currency

func (a ByValueDesc) Len() int           { return len(a) }
func (a ByValueDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValueDesc) Less(i, j int) bool { return a[i].Value > a[j].Value }

func ProcessAndSort(xmlCurrencies []xml.ParsedValute) ([]Currency, error) {
	currencies := make([]Currency, len(xmlCurrencies))

	for index, curr := range xmlCurrencies {
		if curr.NumCode == "" {
			return nil, fmt.Errorf("empty 'NumCode' for currency %s", curr.CharCode)
		}

		numCode, err := strconv.Atoi(curr.NumCode)
		if err != nil {
			return nil, fmt.Errorf("invalid 'NumCode' format '%s': %w", curr.NumCode, err)
		}

		currencies[index] = Currency{
			NumCode:  numCode,
			CharCode: curr.CharCode,
			Value:    curr.Value,
		}
	}

	sort.Sort(ByValueDesc(currencies))

	return currencies, nil
}
