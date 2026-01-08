package main

import (
	"student/currency-processor/internal/conf"
	"student/currency-processor/internal/utils"
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