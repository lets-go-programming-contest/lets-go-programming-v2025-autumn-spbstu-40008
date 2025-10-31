package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  int     `xml:"NumCode" json:"num_code"`
	CharCode string  `xml:"CharCode" json:"char_code"`
	Value    float64 `xml:"Value" json:"value"`
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.Parse()
	if configPath == "" {
		panic("Flag -config is required")
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	xmlDate, err := LoadXMLFile(config.InputFile)
	if err != nil {
		panic(err)
	}

	var valCurs ValCurs
	err = xml.Unmarshal(xmlDate, &valCurs)
	if err != nil {
		panic(fmt.Errorf("unmarshal XML data failed: %w", err))
	}

	SortCurrenciesByValue(valCurs.Valutes)

	err = WriteJSONFile(config.OutputFile, valCurs.Valutes)
	if err != nil {
		panic(err)
	}
}

func LoadConfig(confPath string) (Config, error) {
	data, err := os.ReadFile(confPath)
	if err != nil {
		return Config{}, fmt.Errorf("read yaml file from %q fail cause: %w", confPath, err)
	}
	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return conf, fmt.Errorf("unmarshal yaml file from %q fail: %w", confPath, &err)
	}
	return conf, nil
}

func LoadXMLFile(inputPath string) ([]byte, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("read xml file from %q failed: %w", inputPath, err)
	}
	return data, nil
}

func (val *Valute) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type TempValute struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var tmp TempValute
	if err := decoder.DecodeElement(&tmp, &start); err != nil {
		return fmt.Errorf("failed to decode Valute element: %w", err)
	}

	if tmp.NumCode == "" {
		val.NumCode = 0
	} else {
		numCode, err := strconv.Atoi(tmp.NumCode)
		if err != nil {
			return fmt.Errorf("failed to parse NumCode '%s': %w", tmp.NumCode, err)
		}
		val.NumCode = numCode
	}

	val.CharCode = tmp.CharCode

	if tmp.Value == "" {
		val.Value = 0.0
	} else {
		normValue := strings.ReplaceAll(tmp.Value, ",", ".")
		parsedValue, err := strconv.ParseFloat(normValue, 64)
		if err != nil {
			return fmt.Errorf("failed to parse Value '%s': %w", tmp.Value, err)
		}
		val.Value = parsedValue
	}

	return nil
}

func SortCurrenciesByValue(valutes []Valute) {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})
}

func WriteJSONFile(filePath string, valutes []Valute) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory for %q: %w", filePath, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create JSON file %q: %w", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(valutes); err != nil {
		return fmt.Errorf("encode to JSON file %q: %w", filePath, err)
	}

	return nil
}
