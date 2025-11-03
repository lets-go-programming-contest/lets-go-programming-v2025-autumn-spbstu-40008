package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	filePerm = 0o644
	dirPerm  = 0o755
)

var ErrDataIsNilRefine = errors.New("данные о валютах пусты в refine.go")

// SortByRateValue сортирует список значений по значению, по убыванию
func SortByRateValue(data *CurrencyIndex) error {
	if data == nil {
		return fmt.Errorf("%w", ErrDataIsNilRefine)
	}

	sort.Slice(data.Currencies, func(i, j int) bool {
		return data.Currencies[i].Value > data.Currencies[j].Value
	})

	return nil
}

// SaveToJSON сохраняет отсортированные данные в JSON файл
func SaveToJSON(filePath string, data *CurrencyIndex) error {
	if data == nil {
		return fmt.Errorf("%w", ErrDataIsNilRefine)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("ошибка создания директории: %w", err)
	}

	jsonData, err := json.MarshalIndent(data.Currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка маршалинга в JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, filePerm); err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}