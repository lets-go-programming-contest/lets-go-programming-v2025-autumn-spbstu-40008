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
	Nominal  string `xml:"Nominal"`
	Value    string `xml:"Value"`
}

type CurrencyOutput struct {
	NumCode  int     `json:"num_code"`
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

	valCurs, err := decodeXML(xmlData)
	if err != nil {
		return fmt.Errorf("ошибка декодирования XML: %w", err)
	}

	outputData := make([]CurrencyOutput, 0, len(valCurs.Valutes))

	for _, valute := range valCurs.Valutes {
		// Парсинг Value. Используем "сырое" значение без деления на Номинал.
		value, err := parseValue(valute.Value)
		if err != nil {
			continue
		}

		// Парсинг NumCode.
		valute.NumCode = keepDigits(valute.NumCode)

		numCode, err := strconv.Atoi(valute.NumCode)
		if err != nil {
			continue
		}

		outputData = append(outputData, CurrencyOutput{
			NumCode:  numCode,
			CharCode: valute.CharCode,
			Value:    value, // Сохраняем Value как есть
		})
	}

	// Сортировка по УБЫВАНИЮ поля Value (Raw Value).
	sort.Slice(outputData, func(i, j int) bool {
		return outputData[i].Value > outputData[j].Value
	})

	if err := saveAsJSON(cfg.OutputFile, outputData); err != nil {
		return fmt.Errorf("ошибка сохранения JSON файла: %w", err)
	}

	return nil
}

func decodeXML(data []byte) (*ValCurs, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
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
		return nil, fmt.Errorf("xml decode error: %w", err)
	}

	return &valCurs, nil
}

// parseValue обрабатывает строку с числом, учитывая запятую или точку.
func parseValue(rawInput string) (float64, error) {
	cleaned := cleanString(rawInput)

	if strings.Contains(cleaned, ",") {
		// Формат с запятой (1.234,56). Убираем точки, меняем запятую на точку.
		cleaned = strings.ReplaceAll(cleaned, ".", "")
		cleaned = strings.ReplaceAll(cleaned, ",", ".")
	} else {
		// Формат с точкой (1,234.56). Убираем запятые.
		cleaned = strings.ReplaceAll(cleaned, ",", "")
	}

	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга числа: %w", err)
	}

	return val, nil
}

// keepDigits оставляет в строке только цифры.
func keepDigits(s string) string {
	var builder strings.Builder

	for _, r := range s {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// cleanString убираем пробелы и неразрывные пробелы.
func cleanString(input string) string {
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\t", "")
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\u00A0", "")

	return strings.TrimSpace(input)
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
