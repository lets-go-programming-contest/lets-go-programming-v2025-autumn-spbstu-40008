package main

import (
	"encoding/json"
	"flag"
	"sort"

	"github.com/task-3/internal/decoder"
	"github.com/task-3/internal/output"

	"github.com/task-3/internal/config"
	"github.com/task-3/internal/structures"
)

func sortValuteByValue(val structures.ValCurs) {
	sort.Slice(val.Valute, func(i, j int) bool {
		return val.Valute[i].Value > val.Valute[j].Value
	})
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "", "Path to the YAML configuration file")
	flag.Parse()

	cfg, err := config.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	val, err := decoder.DecodeXML(cfg)
	if err != nil {
		panic(err)
	}

	sortValuteByValue(val)

	jsonData, err := json.MarshalIndent(val.Valute, "", "  ")
	if err != nil {
		panic(err)
	}

	outputFile, err := output.CreateOutputFile(cfg.Output)
	if err != nil {
		panic(err)
	}

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
