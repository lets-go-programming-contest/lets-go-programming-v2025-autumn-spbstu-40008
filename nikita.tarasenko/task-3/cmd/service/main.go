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

func processValutes(valutes []Valute) []Currency {
	currencies := make([]Currency, 0, len(valutes))
	for _, v := range valutes {
		charCode := strings.TrimSpace(v.CharCode)
		valueStr := strings.TrimSpace(v.Value)

		if charCode == "" || valueStr == "" {
			continue
		}

		digits := ""
		for _, r := range v.NumCode {
			if r >= '0' && r <= '9' {
				digits += string(r)
			}
		}

		numCode := 0
		if digits != "" {
			clean := strings.TrimLeft(digits, "0")
			if clean == "" {
				clean = "0"
			}
			if n, err := strconv.Atoi(clean); err == nil {
				numCode = n
			}
		}

		valueStr = strings.ReplaceAll(valueStr, ",", ".")
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		currencies = append(currencies, Currency{
			NumCode:  numCode,
			CharCode: charCode,
			Value:    value,
		})
	}

	if len(currencies) == 0 {
		panic("no valid currencies found in XML")
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies
}

func fixXMLContent(content []byte) []byte {
	content = bytes.Replace(
		content,
		[]byte(`<?xml version="1.0" encoding="windows-1251"?>`),
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>`),
		1,
	)
	content = bytes.ReplaceAll(content, []byte(`encoding="windows-1251"`), []byte(`encoding="UTF-8"`))
	return content
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
	rawContent, err := io.ReadAll(transform.NewReader(xmlFile, decoder))
	if err != nil {
		panic(fmt.Sprintf("failed to decode XML file from windows-1251: %v", err))
	}

	content := fixXMLContent(rawContent)

	var valCurs ValCurs
	if err := xml.Unmarshal(content, &valCurs); err != nil {
		panic(fmt.Sprintf("failed to unmarshal XML: %v", err))
	}

	if valCurs.XMLName.Local != "ValCurs" {
		panic("XML root element is not ValCurs, invalid signature")
	}

	currencies := processValutes(valCurs.Valutes)

	outputDir := filepath.Dir(config.OutputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	outFile, err := os.Create(config.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("failed to create output file: %v", err))
	}
	defer func() { _ = outFile.Close() }()

	encoder := json.NewEncoder(outFile)
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

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
