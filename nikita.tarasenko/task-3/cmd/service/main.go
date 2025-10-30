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

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("flag --config is required")
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	xmlFile, err := os.Open(config.InputFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read input file: %v", err))
	}
	defer xmlFile.Close()

	decoder := charmap.Windows1251.NewDecoder()
	content, err := io.ReadAll(transform.NewReader(xmlFile, decoder))
	if err != nil {
		panic(fmt.Sprintf("failed to decode XML file from windows-1251: %v", err))
	}

	content = bytes.Replace(content, []byte(`<?xml version="1.0" encoding="windows-1251"?>`), []byte(`<?xml version="1.0" encoding="UTF-8"?>`), 1)
	content = bytes.Replace(content, []byte(`encoding="windows-1251"`), []byte(`encoding="UTF-8"`), -1)

	var valCurs ValCurs
	err = xml.Unmarshal(content, &valCurs)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal XML: %v", err))
	}

	if valCurs.XMLName.Local != "ValCurs" {
		panic("XML root element is not ValCurs, invalid signature")
	}

	// Собираем только валидные валюты
	var currencies []Currency
	for _, v := range valCurs.Valutes {
		numCodeStr := strings.TrimSpace(v.NumCode)
		charCodeStr := strings.TrimSpace(v.CharCode)
		valueStr := strings.TrimSpace(v.Value)

		// Пропускаем валюту, если хоть одно поле пустое
		if numCodeStr == "" || charCodeStr == "" || valueStr == "" {
			continue
		}

		cleanNum := strings.TrimLeft(numCodeStr, "0")
		if cleanNum == "" {
			cleanNum = "0"
		}
		numCode, err := strconv.Atoi(cleanNum)
		if err != nil {
			continue
		}

		valueStr = strings.Replace(valueStr, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue // не число — пропускаем
		}

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: charCodeStr,
			Value:    value,
		})
	}

	// После фильтрации: если ни одной валюты нет — ошибка
	if len(currencies) == 0 {
		panic("no valid currencies found in XML")
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

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
