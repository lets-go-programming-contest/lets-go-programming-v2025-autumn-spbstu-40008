package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"io"
	"os"
	"path/filepath"
	"sort"
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

func decodeXML(cfg structures.File) structures.ValCurs {
	xmlFile, err := os.Open(cfg.Input)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := xmlFile.Close(); err != nil {
			return
		}
	}()

	var val structures.ValCurs

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

func sortValuteByValue(val structures.ValCurs) {
	sort.Slice(val.Valute, func(i, j int) bool {
		return val.Valute[i].Value > val.Valute[j].Value
	})
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

	flag.StringVar(&configPath, "config", "", "Path to the YAML configuration file")
	flag.Parse()

	cfg := readFile(configPath)
	val := decodeXML(cfg)

	sortValuteByValue(val)

	jsonData, err := json.MarshalIndent(val.Valute, "", "  ")
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
