package currency

import "sort"

func SortValutes(valutes []*ResultValute) {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})
}
