package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	cpOffset        = 0xC0
	cpA8            = 0xA8
	cpB8            = 0xB8
	asciiLimit      = 0x80
	bufRuneSize     = 4
	maxXMLHeadCheck = 200
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valute  []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  string   `xml:"Nominal"`
	Name     string   `xml:"Name"`
	Value    string   `xml:"Value"`
}

type OutCurrency struct {
	NumCode  int     `json:"num_code" xml:"num_code" yaml:"num_code"`
	CharCode string  `json:"char_code" xml:"char_code" yaml:"char_code"`
	Value    float64 `json:"value" xml:"value" yaml:"value"`
}

type OutCurrenciesXML struct {
	XMLName    xml.Name      `xml:"Currencies"`
	Currencies []OutCurrency `xml:"Currency"`
}

func loadConfig(path string) (Config, error) {
	var cfg Config
	bytes, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("cannot read config file %q: %w", path, err)
	}
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return cfg, fmt.Errorf("cannot parse config yaml: %w", err)
	}
	if cfg.InputFile == "" {
		return cfg, fmt.Errorf("config: input-file is empty")
	}
	if cfg.OutputFile == "" {
		return cfg, fmt.Errorf("config: output-file is empty")
	}
	return cfg, nil
}

func readInputFile(path string) ([]byte, error) {
	inBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read input xml file %q: %w", path, err)
	}
	return inBytes, nil
}

func cp1251ToUTF8(inputBytes []byte) []byte {
	out := make([]byte, 0, len(inputBytes))
	buf := make([]byte, bufRuneSize)
	for idx := range inputBytes {
		b := inputBytes[idx]
		if b < asciiLimit {
			out = append(out, b)
			continue
		}
		var runeVal rune
		switch {
		case b == cpA8:
			runeVal = 0x0401
		case b == cpB8:
			runeVal = 0x0451
		case b >= cpOffset:
			runeVal = 0x0410 + rune(b-cpOffset)
		default:
			runeVal = rune(b)
		}
		n := utf8.EncodeRune(buf, runeVal)
		out = append(out, buf[:n]...)
	}
	return out
}

func detectWindows1251(inputBytes []byte) bool {
	head := strings.ToLower(string(inputBytes))
	if strings.Contains(head, "encoding=\"windows-1251\"") ||
		strings.Contains(head, "encoding='windows-1251'") ||
		strings.Contains(head, "encoding=\"cp1251\"") ||
		strings.Contains(head, "encoding='cp1251'") {
		return true
	}
	if len(inputBytes) > maxXMLHeadCheck {
		head = strings.ToLower(string(inputBytes[:maxXMLHeadCheck]))
		if strings.Contains(head, "encoding=\"windows-1251\"") ||
			strings.Contains(head, "encoding='windows-1251'") {
			return true
		}
	}
	return false
}

func decodeXML(inputBytes []byte) (ValCurs, error) {
	var valcus ValCurs
	if detectWindows1251(inputBytes) {
		inputBytes = cp1251ToUTF8(inputBytes)
	}
	if err := xml.Unmarshal(inputBytes, &valcus); err != nil {
		return valcus, fmt.Errorf("cannot unmarshal input xml: %w", err)
	}
	return valcus, nil
}

func buildOutCurrencies(valcus ValCurs) ([]OutCurrency, error) {
	outList := make([]OutCurrency, 0, len(valcus.Valute))
	for _, val := range valcus.Valute {
		numCodeInt, err := strconv.Atoi(strings.TrimSpace(val.NumCode))
		if err != nil {
			return nil, fmt.Errorf("invalid NumCode %q: %w", val.NumCode, err)
		}

		nominalStr := strings.TrimSpace(val.Nominal)
		if nominalStr == "" {
			nominalStr = "1"
		}

		nominal, err := strconv.Atoi(nominalStr)
		if err != nil || nominal == 0 {
			return nil, fmt.Errorf("invalid Nominal %q for %s: %w", val.Nominal, val.CharCode, err)
		}

		valStr := strings.TrimSpace(val.Value)
		valStr = strings.ReplaceAll(valStr, " ", "")
		valStr = strings.ReplaceAll(valStr, "\u00A0", "")
		valStr = strings.ReplaceAll(valStr, ",", ".")
		parsedValue, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Value %q for %s: %w", val.Value, val.CharCode, err)
		}

		valuePer1 := parsedValue / float64(nominal)

		outCurrency := OutCurrency{
			NumCode:  numCodeInt,
			CharCode: strings.TrimSpace(val.CharCode),
			Value:    valuePer1,
		}
		outList = append(outList, outCurrency)
	}
	sort.Slice(outList, func(i, j int) bool {
		return outList[i].Value > outList[j].Value
	})
	return outList, nil
}

func writeOutput(cfg Config, outList []OutCurrency, format string) error {
	outDir := filepath.Dir(cfg.OutputFile)
	if outDir != "." && outDir != "" {
		if err := os.MkdirAll(outDir, os.FileMode(0o755)); err != nil {
			return fmt.Errorf("cannot create output directory %q: %w", outDir, err)
		}
	}

	outFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		return fmt.Errorf("cannot create output file %q: %w", cfg.OutputFile, err)
	}
	defer func() {
		_ = outFile.Close()
	}()

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		enc := json.NewEncoder(outFile)
		enc.SetIndent("", "  ")
		if err := enc.Encode(outList); err != nil {
			return fmt.Errorf("cannot write json to %q: %w", cfg.OutputFile, err)
		}
	case "yaml", "yml":
		yb, err := yaml.Marshal(outList)
		if err != nil {
			return fmt.Errorf("cannot marshal yaml: %w", err)
		}
		if _, err := outFile.Write(yb); err != nil {
			return fmt.Errorf("cannot write yaml to %q: %w", cfg.OutputFile, err)
		}
	case "xml":
		wrap := OutCurrenciesXML{
			XMLName:    xml.Name{Space: "", Local: "Currencies"},
			Currencies: outList,
		}
		xmlBytes, err := xml.MarshalIndent(wrap, "", "  ")
		if err != nil {
			return fmt.Errorf("cannot marshal xml: %w", err)
		}
		if _, err := io.WriteString(outFile, xml.Header); err != nil {
			return fmt.Errorf("cannot write xml header: %w", err)
		}
		if _, err := outFile.Write(xmlBytes); err != nil {
			return fmt.Errorf("cannot write xml to %q: %w", cfg.OutputFile, err)
		}
	default:
		return fmt.Errorf("unsupported output-format %q (use json, yaml or xml)", format)
	}
	return nil
}

func writeStatus(count int, outPath string, format string) error {
	if _, err := fmt.Fprintf(os.Stdout,
		"Wrote %d records to %s (format=%s)\n",
		count, outPath, format); err != nil {
		return fmt.Errorf("failed to write status to stdout: %w", err)
	}
	return nil
}

func main() {
	cfgPath := flag.String("config", "", "path to YAML config file (contains input-file and output-file)")
	outputFormat := flag.String("output-format", "json", "output format: json (default), yaml or xml")
	flag.Parse()

	if *cfgPath == "" {
		panic("flag -config is required")
	}

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	inBytes, err := readInputFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	valcus, err := decodeXML(inBytes)
	if err != nil {
		panic(err)
	}

	outList, err := buildOutCurrencies(valcus)
	if err != nil {
		panic(err)
	}

	if err := writeOutput(cfg, outList, *outputFormat); err != nil {
		panic(err)
	}

	if err := writeStatus(len(outList), cfg.OutputFile, *outputFormat); err != nil {
		panic(err)
	}
}
