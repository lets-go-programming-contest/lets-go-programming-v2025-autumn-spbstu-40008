package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type Currency struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type ValCurs struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Currencies []Currency `xml:"Valute"`
}

type OutputCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("Config path is required")
	}

	config, err := readConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Error reading config: %v", err))
	}

	currencies, err := parseXML(config.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Error parsing XML: %v", err))
	}

	outputCurrencies := convertAndSortCurrencies(currencies)

	err = saveToJSON(outputCurrencies, config.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("Error saving to JSON: %v", err))
	}

	fmt.Println("Successfully processed currencies")
}

func readConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func parseXML(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	decoder := charmap.Windows1251.NewDecoder()
	utf8Data, err := decoder.Bytes(data)
	if err != nil {
		return nil, err
	}

	xmlContent := string(utf8Data)
	xmlContent = strings.Replace(xmlContent, `encoding="windows-1251"`, `encoding="UTF-8"`, 1)

	var valCurs ValCurs
	err = xml.Unmarshal([]byte(xmlContent), &valCurs)
	if err != nil {
		return nil, fmt.Errorf("xml unmarshal: %w", err)
	}

	return valCurs.Currencies, nil
}

func convertAndSortCurrencies(currencies []Currency) []OutputCurrency {
	output := make([]OutputCurrency, 0, len(currencies))

	for _, currency := range currencies {
		numCode, err := strconv.Atoi(currency.NumCode)
		if err != nil {
			continue
		}

		valueStr := strings.Replace(currency.Value, ",", ".", -1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		output = append(output, OutputCurrency{
			NumCode:  numCode,
			CharCode: currency.CharCode,
			Value:    value,
		})
	}

	for i := 0; i < len(output); i++ {
		for j := i + 1; j < len(output); j++ {
			if output[i].Value < output[j].Value {
				output[i], output[j] = output[j], output[i]
			}
		}
	}

	return output
}

func saveToJSON(currencies []OutputCurrency, path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(currencies)
}
