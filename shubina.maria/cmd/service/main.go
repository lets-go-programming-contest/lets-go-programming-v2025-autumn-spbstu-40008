package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Config для хранения настроек из config.yaml
type Config struct {
	InputFile  string
	OutputFile string
}

// ValCurs для XML структуры от ЦБ
type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

// Valute для одной валюты
type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	ValueStr string `xml:"Value"` // Читаем как строку, потом преобразуем
}

// ValuteJSON для конечного JSON
type ValuteJSON struct {
	NumCode  string  `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	fmt.Println("Запуск программы...")

	// Шаг 1: Читаем конфиг
	config := readConfig("config.yaml")

	// Шаг 2: Читаем и парсим XML
	valutes := parseXML(config.InputFile)

	// Шаг 3: Сортируем по значению (убывание)
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})

	// Шаг 4: Сохраняем в JSON
	saveJSON(config.OutputFile, valutes)

	fmt.Println("Готово! Результат в:", config.OutputFile)
}

func readConfig(filename string) Config {
	// Простой чтение конфига без внешних библиотек
	data, err := os.ReadFile(filename)
	if err != nil {
		panic("Ошибка чтения config.yaml: " + err.Error())
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	var config Config
	for _, line := range lines {
		if strings.Contains(line, "input-file:") {
			config.InputFile = strings.TrimSpace(strings.Split(line, ":")[1])
			config.InputFile = strings.Trim(config.InputFile, " \"")
		}
		if strings.Contains(line, "output-file:") {
			config.OutputFile = strings.TrimSpace(strings.Split(line, ":")[1])
			config.OutputFile = strings.Trim(config.OutputFile, " \"")
		}
	}

	return config
}

func parseXML(filename string) []ValuteJSON {
	// Читаем XML файл
	data, err := os.ReadFile(filename)
	if err != nil {
		panic("Ошибка чтения XML файла: " + err.Error())
	}

	var valCurs ValCurs
	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		panic("Ошибка парсинга XML: " + err.Error())
	}

	// Конвертируем в JSON структуру с правильными типами
	var result []ValuteJSON
	for _, valute := range valCurs.Valutes {
		// Заменяем запятую на точку и конвертируем в float64
		valueStr := strings.Replace(valute.ValueStr, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			fmt.Printf("Ошибка конвертации значения %s: %v\n", valute.ValueStr, err)
			continue
		}

		result = append(result, ValuteJSON{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	return result
}

func saveJSON(filename string, valutes []ValuteJSON) {
	// Создаем папку если нужно
	dir := "result"
	os.MkdirAll(dir, 0755)

	// Создаем файл
	file, err := os.Create(filename)
	if err != nil {
		panic("Ошибка создания JSON файла: " + err.Error())
	}
	defer file.Close()

	// Кодируем в JSON с красивым форматированием
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(valutes)
	if err != nil {
		panic("Ошибка кодирования JSON: " + err.Error())
	}
}
