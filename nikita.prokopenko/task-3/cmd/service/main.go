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
	NumCode  int     `json:"num_code" yaml:"num_code" xml:"num_code"`
	CharCode string  `json:"char_code" yaml:"char_code" xml:"char_code"`
	Value    float64 `json:"value" yaml:"value" xml:"value"`
}

type OutCurrenciesXML struct {
	XMLName    xml.Name      `xml:"Currencies"`
	Currencies []OutCurrency `xml:"Currency"`
}

func cp1251ToUTF8(in []byte) []byte {
	var out []byte
	for i := 0; i < len(in); i++ {
		b := in[i]
		if b < 0x80 {
			out = append(out, b)
			continue
		}
		var r rune
		switch {
		case b == 0xA8:
			r = 0x0401
		case b == 0xB8:
			r = 0x0451
		case b >= 0xC0:
			r = 0x0410 + rune(b-0xC0)
		default:
			r = rune(b)
		}
		buf := make([]byte, 4)
		n := utf8.EncodeRune(buf, r)
		out = append(out, buf[:n]...)
	}
	return out
}

func detectWindows1251(in []byte) bool {
	head := strings.ToLower(string(in))
	if strings.Contains(head, "encoding=\"windows-1251\"") || strings.Contains(head, "encoding='windows-1251'") || strings.Contains(head, "encoding=\"cp1251\"") || strings.Contains(head, "encoding='cp1251'") {
		return true
	}
	if len(in) > 200 {
		head = strings.ToLower(string(in[:200]))
		if strings.Contains(head, "encoding=\"windows-1251\"") || strings.Contains(head, "encoding='windows-1251'") {
			return true
		}
	}
	return false
}

func main() {
	cfgPath := flag.String("config", "", "path to YAML config file (contains input-file and output-file)")
	outputFormat := flag.String("output-format", "json", "output format: json (default), yaml or xml")
	flag.Parse()

	if *cfgPath == "" {
		panic("flag -config is required")
	}

	cfgBytes, err := os.ReadFile(*cfgPath)
	if err != nil {
		panic(fmt.Sprintf("cannot read config file %q: %v", *cfgPath, err))
	}

	var cfg Config
	if err := yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		panic(fmt.Sprintf("cannot parse config yaml: %v", err))
	}
	if cfg.InputFile == "" {
		panic("config: input-file is empty")
	}
	if cfg.OutputFile == "" {
		panic("config: output-file is empty")
	}

	inBytes, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		panic(fmt.Sprintf("cannot read input xml file %q: %v", cfg.InputFile, err))
	}

	if detectWindows1251(inBytes) {
		inBytes = cp1251ToUTF8(inBytes)
	}

	var valcus ValCurs
	if err := xml.Unmarshal(inBytes, &valcus); err != nil {
		panic(fmt.Sprintf("cannot unmarshal input xml: %v", err))
	}

	outCurrencies := make([]OutCurrency, 0, len(valcus.Valute))
	for _, val := range valcus.Valute {
		numCodeInt, errConv := strconv.Atoi(strings.TrimSpace(val.NumCode))
		if errConv != nil {
			panic(fmt.Sprintf("invalid NumCode %q: %v", val.NumCode, errConv))
		}

		nominalStr := strings.TrimSpace(val.Nominal)
		if nominalStr == "" {
			nominalStr = "1"
		}
		nominal, errConv := strconv.Atoi(nominalStr)
		if errConv != nil || nominal == 0 {
			panic(fmt.Sprintf("invalid Nominal %q for %s: %v", val.Nominal, val.CharCode, errConv))
		}

		valStr := strings.TrimSpace(val.Value)
		valStr = strings.ReplaceAll(valStr, " ", "")
		valStr = strings.ReplaceAll(valStr, "\u00A0", "")
		valStr = strings.ReplaceAll(valStr, ",", ".")
		parsedValue, errConv := strconv.ParseFloat(valStr, 64)
		if errConv != nil {
			panic(fmt.Sprintf("invalid Value %q for %s: %v", val.Value, val.CharCode, errConv))
		}

		valuePer1 := parsedValue / float64(nominal)

		outCurrency := OutCurrency{
			NumCode:  numCodeInt,
			CharCode: strings.TrimSpace(val.CharCode),
			Value:    valuePer1,
		}
		outCurrencies = append(outCurrencies, outCurrency)
	}

	sort.Slice(outCurrencies, func(i, j int) bool {
		return outCurrencies[i].Value > outCurrencies[j].Value
	})

	outDir := filepath.Dir(cfg.OutputFile)
	if outDir != "." && outDir != "" {
		if err := os.MkdirAll(outDir, os.FileMode(0755)); err != nil {
			panic(fmt.Sprintf("cannot create output directory %q: %v", outDir, err))
		}
	}

	outFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		panic(fmt.Sprintf("cannot create output file %q: %v", cfg.OutputFile, err))
	}
	defer func() {
		if cerr := outFile.Close(); cerr != nil {
			panic(fmt.Sprintf("error closing output file: %v", cerr))
		}
	}()

	switch strings.ToLower(strings.TrimSpace(*outputFormat)) {
	case "json":
		enc := json.NewEncoder(outFile)
		enc.SetIndent("", "  ")
		if err := enc.Encode(outCurrencies); err != nil {
			panic(fmt.Sprintf("cannot write json to %q: %v", cfg.OutputFile, err))
		}
	case "yaml", "yml":
		yb, err := yaml.Marshal(outCurrencies)
		if err != nil {
			panic(fmt.Sprintf("cannot marshal yaml: %v", err))
		}
		if _, err := outFile.Write(yb); err != nil {
			panic(fmt.Sprintf("cannot write yaml to %q: %v", cfg.OutputFile, err))
		}
	case "xml":
		wrap := OutCurrenciesXML{
			XMLName:    xml.Name{Local: "Currencies"},
			Currencies: outCurrencies,
		}
		xmlBytes, err := xml.MarshalIndent(wrap, "", "  ")
		if err != nil {
			panic(fmt.Sprintf("cannot marshal xml: %v", err))
		}
		if _, err := io.WriteString(outFile, xml.Header); err != nil {
			panic(fmt.Sprintf("cannot write xml header: %v", err))
		}
		if _, err := outFile.Write(xmlBytes); err != nil {
			panic(fmt.Sprintf("cannot write xml to %q: %v", cfg.OutputFile, err))
		}
	default:
		panic(fmt.Sprintf("unsupported output-format %q (use json, yaml or xml)", *outputFormat))
	}

	if _, err := fmt.Fprintf(os.Stdout, "Wrote %d records to %s (format=%s)\n", len(outCurrencies), cfg.OutputFile, *outputFormat); err != nil {
		panic(fmt.Sprintf("failed to write status to stdout: %v", err))
	}
}
