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

func processValutes(inputValutes []Valute) []Currency {
	currencies := make([]Currency, 0, len(inputValutes))

	for _, valuteItem := range inputValutes {
		charCode := strings.TrimSpace(valuteItem.CharCode)

		digits := ""

		for _, ch := range valuteItem.NumCode {
			if ch >= '0' && ch <= '9' {
				digits += string(ch)
			}
		}

		numCode := 0

		if digits != "" {
			cleaned := strings.TrimLeft(digits, "0")
			if cleaned == "" {
				cleaned = "0"
			}

			parsedNum, err := strconv.Atoi(cleaned)
			if err == nil {
				numCode = parsedNum
			}
		}

		valueStr := strings.TrimSpace(valuteItem.Value)
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

func prepareXMLContent(data []byte) []byte {
	data = bytes.Replace(
		data,
		[]byte(`<?xml version="1.0" encoding="windows-1251"?>`),
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>`),
		1,
	)
	data = bytes.ReplaceAll(data, []byte(`encoding="windows-1251"`), []byte(`encoding="UTF-8"`))

	return data
}

func main() {
	cfgPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *cfgPath == "" {
		panic("flag --config is required")
	}

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	xmlFile, err := os.Open(cfg.InputFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read input file: %v", err))
	}

	defer func() { _ = xmlFile.Close() }()

	decoder := charmap.Windows1251.NewDecoder()
	rawData, err := io.ReadAll(transform.NewReader(xmlFile, decoder))
	if err != nil {
		panic(fmt.Sprintf("failed to decode XML file from windows-1251: %v", err))
	}

	cleanData := prepareXMLContent(rawData)

	var valCurs ValCurs
	if err := xml.Unmarshal(cleanData, &valCurs); err != nil {
		panic(fmt.Sprintf("failed to unmarshal XML: %v", err))
	}

	if valCurs.XMLName.Local != "ValCurs" {
		panic("XML root element is not ValCurs, invalid signature")
	}

	currencies := processValutes(valCurs.Valutes)

	const dirPerm = 0o755

	outputDir := filepath.Dir(cfg.OutputFile)
	if err := os.MkdirAll(outputDir, dirPerm); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	outFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("failed to create output file: %v", err))
	}

	defer func() { _ = outFile.Close() }()

	jsonEncoder := json.NewEncoder(outFile)
	jsonEncoder.SetIndent("", "  ")

	if err := jsonEncoder.Encode(currencies); err != nil {
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
