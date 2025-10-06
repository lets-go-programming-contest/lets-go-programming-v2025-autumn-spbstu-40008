package main

import "fmt"

const MinTempDefault = 15
const MaxTempDefault = 30
const ErrVal = -1

func main() {
	var countDepartments int
	var countEmployees int
	var temp int
	var newMinTemperature int
	var newMaxTemperature int
	var setTempSign string

	fmt.Scan(&countDepartments)
	for range countDepartments {
		
		fmt.Scan(&countEmployees)
		for range countEmployees {
			fmt.Scan(&setTempSign, &temp)
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
			fmt.Println(newMinTemperature)
		}
	}
}