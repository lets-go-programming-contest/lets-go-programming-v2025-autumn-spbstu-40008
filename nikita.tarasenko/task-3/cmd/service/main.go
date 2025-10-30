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

func processValute(v Valute) (*Currency, bool) {
	charCodeStr := strings.TrimSpace(v.CharCode)
	valueStr := strings.TrimSpace(v.Value)

	if charCodeStr == "" || valueStr == "" {
		return nil, false
	}

	// Извлекаем только цифры из NumCode
	digitsOnly := ""
	for _, r := range v.NumCode {
		if r >= '0' && r <= '9' {
			digitsOnly += string(r)
		}
	}

	var numCode int
	if digitsOnly == "" {
		numCode = 0
	} else {
		cleanNum := strings.TrimLeft(digitsOnly, "0")
		if cleanNum == "" {
			cleanNum = "0"
		}
		if n, err := strconv.Atoi(cleanNum); err == nil {
			numCode = n
		} else {
			numCode = 0
		}
	}

	valueStr = strings.ReplaceAll(valueStr, ",", ".")
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, false
	}

	return &Currency{
		NumCode:  numCode,
		CharCode: charCodeStr,
		Value:    value,
	}, true
}

func parseXML(content []byte) ([]Currency, error) {
	var valCurs ValCurs
	if err := xml.Unmarshal(content, &valCurs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	if valCurs.XMLName.Local != "ValCurs" {
		return nil, fmt.Errorf("XML root element is not ValCurs, invalid signature")
	}

	currencies := make([]Currency, 0, len(valCurs.Valutes))
	for _, v := range valCurs.Valutes {
		if cur, ok := processValute(v); ok {
			currencies = append(currencies, *cur)
		}
	}

	if len(currencies) == 0 {
		return nil, fmt.Errorf("no valid currencies found in XML")
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies, nil
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
	defer func() { _ = xmlFile.Close() }()

	decoder := charmap.Windows1251.NewDecoder()
	content, err := io.ReadAll(transform.NewReader(xmlFile, decoder))
	if err != nil {
		panic(fmt.Sprintf("failed to decode XML file from windows-1251: %v", err))
	}

	// Исправляем объявление кодировки
	content = bytes.ReplaceAll(
		content,
		[]byte(`<?xml version="1.0" encoding="windows-1251"?>`),
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>`),
	)
	content = bytes.ReplaceAll(
		content,
		[]byte(`encoding="windows-1251"`),
		[]byte(`encoding="UTF-8"`),
	)

	currencies, err := parseXML(content)
	if err != nil {
		panic(err)
	}

	outputDir := filepath.Dir(config.OutputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	file, err := os.Create(config.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("failed to create output file: %v", err))
	}
	defer func() { _ = file.Close() }()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(currencies); err != nil {
		panic(fmt.Sprintf("failed to encode JSON: %v", err))
	}
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
