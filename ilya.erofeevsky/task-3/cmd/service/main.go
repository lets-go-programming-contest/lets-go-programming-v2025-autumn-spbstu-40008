package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v2"

	"task-3/internal/structures"
)

var (
	ErrUnsupportedCharset = errors.New("unsupported charset")
)

func ReadFile(configPath string) structures.File {
	var cfg structures.File

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Error reading config file %s: %v", configPath, err))
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Error decoding YAML config: %v", err))
	}

	if cfg.Input == "" || cfg.Output == "" {
		panic("Config file must contain 'input-file' and 'output-file' paths.")
	}

	return cfg
}

func decodeXML(cfg structures.File) structures.ReadingXML {
	xmlFile, err := os.Open(cfg.Input)
	if err != nil {
		panic(fmt.Sprintf("Error opening XML input file %s: %v", cfg.Input, err))
	}

	defer func() {
		if err := xmlFile.Close(); err != nil {
			fmt.Printf("Error closing XML file: %v\n", err)
		}
	}()

	var xmlData structures.ReadingXML

	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, ErrUnsupportedCharset
	}

	err = decoder.Decode(&xmlData)
	if err != nil {
		panic(fmt.Sprintf("Error decoding XML from file %s: %v", cfg.Input, err))
	}

	return xmlData
}

func SortAndProcessCurrencies(xmlData structures.ReadingXML) []structures.ProcessedCurrency {
	processed := make([]structures.ProcessedCurrency, 0, len(xmlData.Information))

	for _, item := range xmlData.Information {
		stringValue := strings.ReplaceAll(item.Value, ",", ".")
		value, errValue := strconv.ParseFloat(stringValue, 64)

		if errValue != nil {
			continue
		}

		numCode := 0
		if item.NumCode != "" {
			parsed, errNumCode := strconv.Atoi(strings.TrimSpace(item.NumCode))
			if errNumCode == nil {
				numCode = parsed
			}
		}

		processed = append(processed, structures.ProcessedCurrency{
			NumCode:  numCode,
			CharCode: strings.TrimSpace(item.CharCode),
			Value:    value,
		})
	}

	sort.Slice(processed, func(i, j int) bool {
		return processed[i].Value > processed[j].Value
	})

	return processed
}

func createOutputFile(filename string) *os.File {
	dirPath := filepath.Dir(filename)

	const DirPerm = 0o755

	if err := os.MkdirAll(dirPath, DirPerm); err != nil {
		panic(fmt.Sprintf("Error creating output directory %s: %v", dirPath, err))
	}

	const FilePerm = 0o644

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePerm)
	if err != nil {
		panic(fmt.Sprintf("Error opening/creating output file %s: %v", filename, err))
	}

	return file
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "Path to the YAML configuration file")
	flag.Parse()

	if configPath == "" {
		panic("Flag --config must be set to the path of the YAML configuration file.")
	}

	cfg := ReadFile(configPath)
	xmlData := decodeXML(cfg)
	sortedCurrencies := SortAndProcessCurrencies(xmlData)

	resultItems := make([]structures.ResultItem, 0, len(sortedCurrencies))
	for _, curr := range sortedCurrencies {
		resultItems = append(resultItems, structures.ResultItem{
			NumCode:  curr.NumCode,
			CharCode: curr.CharCode,
			Value:    curr.Value,
		})
	}

	jsonData, err := json.MarshalIndent(resultItems, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("Error marshaling data to JSON: %v", err))
	}

	outputFile := createOutputFile(cfg.Output)
	defer func() {
		if err := outputFile.Close(); err != nil {
			fmt.Printf("Error closing output file: %v\n", err)
		}
	}()

	if _, err := outputFile.Write(jsonData); err != nil {
		panic(fmt.Sprintf("Error writing JSON data to file %s: %v", cfg.Output, err))
	}
}
