package app

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

var ErrUnsupportedCharset = errors.New("неподдерживаемая кодировка")

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

	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedCharset, charset)
		}
	}

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return fmt.Errorf("ошибка декодирования XML: %w", err)
	}

	outputData := make([]CurrencyOutput, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		valueStr := strings.Replace(valute.Value, ",", ".", 1)

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}

		outputData = append(outputData, CurrencyOutput{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value / float64(valute.Nominal),
		})
	}

	sort.Slice(outputData, func(i, j int) bool {
		return outputData[i].Value > outputData[j].Value
	})

	if err := saveAsJSON(cfg.OutputFile, outputData); err != nil {
		return fmt.Errorf("ошибка сохранения JSON файла: %w", err)
	}

	return nil
}

func saveAsJSON(path string, data []CurrencyOutput) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("ошибка создания директории: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("ошибка кодирования JSON: %w", err)
	}

	return nil
}
