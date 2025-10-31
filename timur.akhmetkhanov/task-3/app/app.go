package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	Nominal  int    `xml:"Nominal"`
	Value    string `xml:"Value"`
}

type CurrencyOutput struct {
	NumCode  string  `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	return &cfg, nil
}

func Run(cfg *Config) error {
	xmlData, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения XML файла: %w", err)
	}

	var valCurs ValCurs
	if err := xml.Unmarshal(xmlData, &valCurs); err != nil {
		return fmt.Errorf("ошибка декодирования XML: %w", err)
	}

	outputData := make([]CurrencyOutput, 0, len(valCurs.Valutes))
	for _, v := range valCurs.Valutes {
		valueStr := strings.Replace(v.Value, ",", ".", 1)
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		outputData = append(outputData, CurrencyOutput{
			NumCode:  v.NumCode,
			CharCode: v.CharCode,
			Value:    value / float64(v.Nominal),
		})
	}

	sort.Slice(outputData, func(i, j int) bool {
		return outputData[i].Value > outputData[j].Value
	})

	if err := saveAsJSON(cfg.OutputFile, outputData); err != nil {
		return fmt.Errorf("ошибка сохранения JSON файла: %w", err)
	}

	fmt.Printf("Результат успешно сохранен в %s\n", cfg.OutputFile)
	return nil
}

func saveAsJSON(path string, data []CurrencyOutput) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
