package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/task-3/internal/structures"

	"golang.org/x/text/encoding/charmap"

	"gopkg.in/yaml.v2"
)

func readFile(configPath string) structures.File {
	var cfg structures.File

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func decodeXML(cfg structures.File) structures.ValCursXML {
	xmlFile, err := os.Open(cfg.Input)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := xmlFile.Close(); err != nil {
			return
		}
	}()

	var val structures.ValCursXML

	decoder := xml.NewDecoder(xmlFile)

	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}

		return nil, io.EOF
	}

	err = decoder.Decode(&val)
	if err != nil {
		panic(err)
	}

	return val
}

func normalizeValues(val structures.ValCursXML) structures.ValCursJSON {
	jsonValCurs := structures.ValCursJSON{
		Valute: make([]structures.ValuteJSON, 0, len(val.Valute)),
	}

	for _, xmlValute := range val.Valute {
		cleanValueStr := strings.ReplaceAll(xmlValute.Value, ",", ".")

		valueFloat, err := strconv.ParseFloat(cleanValueStr, 64)
		if err != nil {
			panic(err)
		}

		jsonValCurs.Valute = append(jsonValCurs.Valute, structures.ValuteJSON{
			NumCode:  xmlValute.NumCode,
			CharCode: xmlValute.CharCode,
			Value:    valueFloat,
		})
	}

	return jsonValCurs
}

func sortValuteByValue(val structures.ValCursXML) structures.ValCursJSON {
	valJSON := normalizeValues(val)
	sort.Slice(valJSON.Valute, func(i, j int) bool {
		return valJSON.Valute[i].Value > valJSON.Valute[j].Value
	})

	return valJSON
}

func createOutputFile(filename string) *os.File {
	dirPath := filepath.Dir(filename)

	const DirPerm = 0o755

	if err := os.MkdirAll(dirPath, DirPerm); err != nil {
		panic(err)
	}

	const FilePerm = 0o644

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePerm)
	if err != nil {
		panic(err)
	}

	return file
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "Path to the YAML configuration file")
	flag.Parse()

	cfg := readFile(configPath)
	val := decodeXML(cfg)
	valJSON := sortValuteByValue(val)

	jsonData, err := json.MarshalIndent(valJSON.Valute, "", "  ")
	if err != nil {
		panic(err)
	}

	outputFile := createOutputFile(cfg.Output)
	defer func() {
		if err := outputFile.Close(); err != nil {
			return
		}
	}()

	_, err = outputFile.Write(jsonData)
	if err != nil {
		panic(err)
	}
}
