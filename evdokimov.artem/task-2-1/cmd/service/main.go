package main

import "fmt"

func main() {
	var countDepartments, countWorkers uint

	_, err := fmt.Scan(&countDepartments)
	if err != nil {
		fmt.Println("Invalid departments count")
		return
	}

	var (
		minValue, maxValue uint8
		currentTemp        uint8
		operationSign      string
	)

	for range countDepartments {
		_, err = fmt.Scan(&countWorkers)
		if err != nil {
			fmt.Println("Invalid workers count")
			return
		}

		minValue, maxValue = 15, 30

		for range countWorkers {
			_, err = fmt.Scan(&operationSign, &currentTemp)
			if err != nil {
				fmt.Println("Invalid temperature input")
				return
			}

			switch operationSign {
			case ">=":
				if currentTemp > minValue {
					minValue = currentTemp
				}
			case "<=":
				if currentTemp < maxValue {
					maxValue = currentTemp
				}
			default:
				fmt.Println("Unknown operation")
				return
			}

			if minValue > maxValue {
				fmt.Println(-1)
			} else {
				fmt.Println(minValue)
			}
		}
	}
}
