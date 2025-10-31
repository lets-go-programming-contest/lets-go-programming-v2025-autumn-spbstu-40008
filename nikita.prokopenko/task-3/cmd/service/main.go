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
	XMLName    xml.Name       `xml:"Currencies"`
	Currencies []OutCurrency `xml:"Currency"`
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
		panic(fmt.Errorf("cannot read config file %q: %w", *cfgPath, err))
	}
	var cfg Config
	if err := yaml.Unmarshal(cfgBytes, &cfg); err != nil {
		panic(fmt.Errorf("cannot parse config yaml: %w", err))
	}
	if cfg.InputFile == "" {
		panic("config: input-file is empty")
	}
	if cfg.OutputFile == "" {
		panic("config: output-file is empty")
	}

	inBytes, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		panic(fmt.Errorf("cannot read input xml file %q: %w", cfg.InputFile, err))
	}

	var valcus ValCurs
	if err := xml.Unmarshal(inBytes, &valcus); err != nil {
		panic(fmt.Errorf("cannot unmarshal input xml: %w", err))
	}

	outs := make([]OutCurrency, 0, len(valcus.Valute))
	for _, v := range valcus.Valute {
		numCode, err := strconv.Atoi(strings.TrimSpace(v.NumCode))
		if err != nil {
			panic(fmt.Errorf("invalid NumCode %q: %w", v.NumCode, err))
		}

		nominalStr := strings.TrimSpace(v.Nominal)
		if nominalStr == "" {
			nominalStr = "1"
		}
		nominal, err := strconv.Atoi(nominalStr)
		if err != nil || nominal == 0 {
			panic(fmt.Errorf("invalid Nominal %q for %s: %w", v.Nominal, v.CharCode, err))
		}

		valStr := strings.TrimSpace(v.Value)
		valStr = strings.ReplaceAll(valStr, " ", "")
		valStr = strings.ReplaceAll(valStr, "\u00A0", "")
		valStr = strings.ReplaceAll(valStr, ",", ".")
		value, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			panic(fmt.Errorf("invalid Value %q for %s: %w", v.Value, v.CharCode, err))
		}

		valuePer1 := value / float64(nominal)

		out := OutCurrency{
			NumCode:  numCode,
			CharCode: strings.TrimSpace(v.CharCode),
			Value:    valuePer1,
		}
		outs = append(outs, out)
	}

	sort.Slice(outs, func(i, j int) bool {
		return outs[i].Value > outs[j].Value
	})

	outDir := filepath.Dir(cfg.OutputFile)
	if outDir != "." && outDir != "" {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			panic(fmt.Errorf("cannot create output directory %q: %w", outDir, err))
		}
	}

	f, err := os.Create(cfg.OutputFile)
	if err != nil {
		panic(fmt.Errorf("cannot create output file %q: %w", cfg.OutputFile, err))
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			panic(fmt.Errorf("error closing output file: %w", cerr))
		}
	}()

	switch strings.ToLower(strings.TrimSpace(*outputFormat)) {
	case "json":
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(outs); err != nil {
			panic(fmt.Errorf("cannot write json to %q: %w", cfg.OutputFile, err))
		}
	case "yaml", "yml":
		yb, err := yaml.Marshal(outs)
		if err != nil {
			panic(fmt.Errorf("cannot marshal yaml: %w", err))
		}
		if _, err := f.Write(yb); err != nil {
			panic(fmt.Errorf("cannot write yaml to %q: %w", cfg.OutputFile, err))
		}
	case "xml":
		wrap := OutCurrenciesXML{Currencies: outs}
		xb, err := xml.MarshalIndent(wrap, "", "  ")
		if err != nil {
			panic(fmt.Errorf("cannot marshal xml: %w", err))
		}
		if _, err := io.WriteString(f, xml.Header); err != nil {
			panic(fmt.Errorf("cannot write xml header: %w", err))
		}
		if _, err := f.Write(xb); err != nil {
			panic(fmt.Errorf("cannot write xml to %q: %w", cfg.OutputFile, err))
		}
	default:
		panic(fmt.Errorf("unsupported output-format %q (use json, yaml or xml)", *outputFormat))
	}

	fmt.Fprintf(os.Stdout, "Wrote %d records to %s (format=%s)\n", len(outs), cfg.OutputFile, *outputFormat)
}
