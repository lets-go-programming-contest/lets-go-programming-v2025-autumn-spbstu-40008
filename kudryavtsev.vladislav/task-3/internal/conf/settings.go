package conf

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppSettings struct {
	SourcePath      string `yaml:"input-file"`
	DestinationPath string `yaml:"output-file"`
}

func FetchPathFromArgs() string {
	var flagValue string
	flag.StringVar(&flagValue, "config", "", "path to configuration file")
	flag.Parse()

	return flagValue
}

func LoadSettings(filename string) (AppSettings, error) {
	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		return AppSettings{}, fmt.Errorf("loading config failed: %w", err)
	}

	var settings AppSettings
	err = yaml.Unmarshal(rawBytes, &settings)
	if err != nil {
		return AppSettings{}, fmt.Errorf("parsing yaml failed: %w", err)
	}

	return settings, nil
}