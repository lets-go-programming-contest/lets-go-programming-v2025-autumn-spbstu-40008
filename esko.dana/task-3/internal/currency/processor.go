package currency

import (
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
	var currencies []Currency

	for _, curr := range xmlCurrencies {
		numCode, err := strconv.Atoi(curr.NumCode)
		if err != nil {
			continue
		}

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: curr.CharCode,
			Value:    curr.Value,
		})
	}

	sort.Sort(ByValueDesc(currencies))

	return currencies, nil
}
