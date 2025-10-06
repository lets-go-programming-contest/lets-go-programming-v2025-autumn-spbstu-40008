package main

import "fmt"

const MinTempDefault = 15
const MaxTempDefault = 30
const ErrVal         = -1

func main() {
	var countDepartments  int
	var countEmployees    int
	var temp              int
	var newMinTemperature int
	var newMaxTemperature int
	var setTempSign       string

	if _, err := fmt.Scan(&countDepartments); err != nil {
		return
	}
	for range countDepartments {
		if _, err := fmt.Scan(&countEmployees); err != nil {
			return
		}

		for range countEmployees {
			if _, err := fmt.Scan(&setTempSign, &temp); err != nil {
				return
			}
			switch setTempSign {
			case ">=":
				newMinTemperature = max(MaxTempDefault, temp)
			case "<=":
				newMaxTemperature = min(MinTempDefault, temp)
			default:
				continue
			}
			if newMinTemperature > newMaxTemperature {
				fmt.Println(ErrVal)
			}
			fmt.Println(newTemperature)
		}
	}
}