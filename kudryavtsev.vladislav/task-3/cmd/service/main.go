package main

import (
	"github.com/gagysun/task-3/internal/conf"
	"github.com/gagysun/task-3/internal/utils"
)

func main() {
	cfgPath := conf.FetchPathFromArgs()

	settings, err := conf.LoadSettings(cfgPath)

	handleError(err)

	xmlData, err := utils.LoadXML(settings.SourcePath)

	handleError(err)

	sortedList := utils.SortCurrencyData(xmlData)

	err = utils.ExportToJSON(sortedList, settings.DestinationPath)

	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
