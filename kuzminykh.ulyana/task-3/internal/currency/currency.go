package currency

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kuzminykh.ulyana/task-3/internal/models"
)

func Sort(valutes []models.Valute) ([]models.Output, error) {
	outdata := make([]models.Output, 0, len(valutes))

	for _, valute := range valutes {
		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			return nil, fmt.Errorf("converting num code '%s': %w", valute.NumCode, err)
		}

		valueStr := strings.Replace(valute.Value, ",", ".", 1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing value '%s': %w", valute.Value, err)
		}

		outdata = append(outdata, models.Output{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(outdata, func(i, j int) bool {
		return outdata[i].Value > outdata[j].Value
	})

	return outdata, nil
}
