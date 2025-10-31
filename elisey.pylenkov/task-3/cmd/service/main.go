package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"task-3/internal/config"
	"task-3/internal/parser"
	"task-3/internal/structures"
)

type ByDescendingValue []structures.Valute

func (a ByDescendingValue) Len() int      { return len(a) }
func (a ByDescendingValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDescendingValue) Less(i, j int) bool {
	return float64(a[i].Value) > float64(a[j].Value)
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "config path")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		panic("config error: " + err.Error())
	}

	valCurs, err := parser.ParseCurrencyXML(cfg.InputFile)
	if err != nil {
		panic("xml error: " + err.Error())
	}

	sort.Sort(ByDescendingValue(valCurs.Valutes))

	outputData := parser.ConvertToOutput(valCurs.Valutes)

	jsonData, err := json.MarshalIndent(outputData, "", "	")
	if err != nil {
		panic(fmt.Sprintf("json saving error: %v", err))
	}

	dir := filepath.Dir(cfg.OutputFile)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(fmt.Sprintf("creating output folder error: %v", err))
	}

	err = os.WriteFile(cfg.OutputFile, jsonData, 0644)
	if err != nil {
		panic(fmt.Sprintf("writing output file error: %v", err))
	}
}
