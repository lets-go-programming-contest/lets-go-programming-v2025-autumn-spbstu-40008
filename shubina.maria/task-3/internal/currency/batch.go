package currency

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

// LoadRatesFromFile читает XML файл с курсами валют
func LoadRatesFromFile(filePath string) (*CurrencyIndex, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	decoder := charmap.Windows1251.NewDecoder()
	decodedData, err := decoder.Bytes(data)
	if err != nil {
		return nil, fmt.Errorf("ошибка кодировки: %w", err)
	}

	xmlStr := string(decodedData)
	if idx := strings.Index(xmlStr, "?>"); idx != -1 {
		xmlStr = xmlStr[idx+2:] // Убираем заголовок XML, если есть
	}

	var result CurrencyIndex
	err = xml.Unmarshal([]byte(xmlStr), &result)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга XML: %w", err)
	}

	return &result, nil
}