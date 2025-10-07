package main

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"task-3/structures"
)

func readFile(cfg structures.File) structures.File {
	yamlFile, err := os.ReadFile("config.yaml")
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
	xmlFile, err := os.ReadFile(cfg.Input)
	if err != nil {
		panic(err)
	}

	var val structures.ValCurs

	err = xml.Unmarshal(xmlFile, &val)
	if err != nil {
		panic(err)
	}

	return val
}

func normalizeValues(val []structures.Valute) {
	for i := range val {
		val[i].Value = strings.ReplaceAll(val[i].Value, ",", ".")
	}
}

func sortValuteByValue(val structures.ValCurs) {
	normalizeValues(val.Valute)
	sort.Slice(val.Valute, func(i, j int) bool {

		valI, errI := strconv.ParseFloat(val.Valute[i].Value, 64)
		valJ, errJ := strconv.ParseFloat(val.Valute[j].Value, 64)

		if errI != nil || errJ != nil {
			return false
		}

		return valI > valJ
	})
}

func createOutputFile(filename string) *os.File {
	dirPath := filepath.Dir(filename)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	return file
}

func main() {
	var cfg structures.File
	cfg = readFile(cfg)
	val := decodeXML(cfg)
	sortValuteByValue(val)

	jsonData, err := json.MarshalIndent(val.Valute, "", "  ")

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
