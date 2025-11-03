package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Maria-Mariaa/task-3/internal/config"
	"github.com/Maria-Mariaa/task-3/internal/currency"
)


func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Путь к файлу конфигурации YAML")
	flag.Parse()

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "Ошибка: Не указан флаг --config\n")
		os.Exit(1)
	}

	cfg, err := config.ReadSettings(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: Не удалось загрузить конфиг: %v\n", err)
		os.Exit(1)
	}

	currencyData, err := currency.LoadRatesFromFile(cfg.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: Не удалось прочитать курсы валют: %v\n", err)
		os.Exit(1)
	}

	err = currency.SortByRateValue(currencyData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: Не удалось отсортировать: %v\n", err)
		os.Exit(1)
	}

	err = currency.SaveToJSON(cfg.OutputFile, currencyData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: Не удалось сохранить JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Успешно обработано %d валют. Результат: %s\n",
		len(currencyData.Currencies), cfg.OutputFile)
}
