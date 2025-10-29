package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"gopkg.in/yaml.v3"
)

// Config структура для конфигурации
type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

// ValCurs XML-структура для корневого элемента
type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

// Valute XML-структура для валюты
type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

// Currency структура для результата (JSON)
type Currency struct {
	NumCode  string  `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("flag --config is required")
	}

	// Шаг 1: Чтение конфига
	config, err := loadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	// Шаг 2: Чтение и декодирование XML
	xmlFile, err := os.Open(config.InputFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read input file: %v", err))
	}
	defer xmlFile.Close()

	// Декодируем из windows-1251 в UTF-8
	decoder := charmap.Windows1251.NewDecoder()
	content, err := io.ReadAll(transform.NewReader(xmlFile, decoder))
	if err != nil {
		panic(fmt.Sprintf("failed to decode XML file from windows-1251: %v", err))
	}

	// Удаляем декларацию encoding="windows-1251" из XML
	content = bytes.Replace(content, []byte(`<?xml version="1.0" encoding="windows-1251"?>`), []byte(`<?xml version="1.0" encoding="UTF-8"?>`), 1)
	content = bytes.Replace(content, []byte(`encoding="windows-1251"`), []byte(`encoding="UTF-8"`), -1)

	var valCurs ValCurs
	err = xml.Unmarshal(content, &valCurs)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal XML: %v", err))
	}

	// Проверка сигнатуры данных: наличие корректного корневого элемента и валют
	if valCurs.XMLName.Local != "ValCurs" {
		panic("XML root element is not ValCurs, invalid signature")
	}
	if len(valCurs.Valutes) == 0 {
		panic("XML contains no Valute elements, invalid signature")
	}

	// Шаг 3: Преобразование и сортировка
	currencies := make([]Currency, 0, len(valCurs.Valutes))
	for _, v := range valCurs.Valutes {
		if v.NumCode == "" || v.CharCode == "" || v.Value == "" {
			panic(fmt.Sprintf("Valute with ID %s has missing required fields, invalid signature", v.ID))
		}
		valueStr := strings.TrimSpace(v.Value)
		valueStr = strings.Replace(valueStr, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			panic(fmt.Sprintf("failed to parse value for currency %s: %v", v.CharCode, err))
		}
		currencies = append(currencies, Currency{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    value,
		})
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	// Шаг 4: Запись JSON
	outputDir := filepath.Dir(config.OutputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	file, err := os.Create(config.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("failed to create output file: %v", err))
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(currencies)
	if err != nil {
		panic(fmt.Sprintf("failed to encode JSON: %v", err))
	}
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
