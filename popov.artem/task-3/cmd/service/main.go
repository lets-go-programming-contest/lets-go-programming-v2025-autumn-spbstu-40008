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
	"gopkg.in/yaml.v3"
)

// Configuration holds input and output file paths.
type Configuration struct {
	InputPath  string `yaml:"input-file"`
	OutputPath string `yaml:"output-file"`
}

// LoadConfiguration loads config from YAML file.
func LoadConfiguration(path string) (*Configuration, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Configuration
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &cfg, nil
}

// Float64Custom is a custom float type for XML parsing.
type Float64Custom float64

// CurrencyInfo represents currency data structure.
type CurrencyInfo struct {
	XMLName  xml.Name      `json:"-"         xml:"Valute"`
	ID       string        `json:"-"         xml:"ID,attr"`
	CodeNum  int           `json:"num_code"  xml:"NumCode"`
	CodeChar string        `json:"char_code" xml:"CharCode"`
	Nominal  int           `json:"-"         xml:"Nominal"`
	FullName string        `json:"-"         xml:"Name"`
	Value    Float64Custom `json:"value"     xml:"Value"`
}

// ExchangeRates is a container for currency list.
type ExchangeRates struct {
	XMLName xml.Name       `xml:"ValCurs"`
	Rates   []CurrencyInfo `xml:"Valute"`
}

// UnmarshalXML parses float from XML, handling comma as decimal separator.
func (f *Float64Custom) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	if err := d.DecodeElement(&value, &start); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	value = strings.ReplaceAll(value, ",", ".")
	
	parsedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("float parse error: %w", err)
	}

	*f = Float64Custom(parsedValue)

	return nil
}

// xmlCharsetReader handles different charsets (like windows-1251).
func xmlCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "windows-1251":
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	default:
		return input, nil
	}
}

// ParseXMLData decodes XML data.
func ParseXMLData(data []byte, target interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = xmlCharsetReader

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("XML decode failed: %w", err)
	}

	return nil
}

// WriteJSONToFile saves currencies to JSON file.
func WriteJSONToFile(path string, data []CurrencyInfo) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("create dir failed: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file failed: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode JSON failed: %w", err)
	}

	return nil
}

// ExecuteProcessing processes the XML and saves as JSON.
func ExecuteProcessing(config *Configuration) error {
	xmlContent, err := os.ReadFile(config.InputPath)
	if err != nil {
		return fmt.Errorf("read XML failed: %w", err)
	}

	var rates ExchangeRates
	if err := ParseXMLData(xmlContent, &rates); err != nil {
		return fmt.Errorf("parse XML failed: %w", err)
	}

	currencies := rates.Rates
	sort.Slice(currencies, func(i, j int) bool {
		return float64(currencies[i].Value) > float64(currencies[j].Value)
	})

	if err := WriteJSONToFile(config.OutputPath, currencies); err != nil {
		return fmt.Errorf("write JSON failed: %w", err)
	}

	return nil
}

func main() {
	configFilePath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()

	cfg, err := LoadConfiguration(*configFilePath)
	if err != nil {
		panic(err)
	}

	if err := ExecuteProcessing(cfg); err != nil {
		panic(err)
	}
}
