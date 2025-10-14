package main

import "fmt"

func main() {
	var (
		numberOfDepartments, numberOfEmployees, maxTemp, minTemp, temp int
		data                                                           string
	)

	_, err := fmt.Scan(&numberOfDepartments)
	if err != nil {
		return
	}

	for range numberOfDepartments {
		_, err := fmt.Scan(&numberOfEmployees)
		if err != nil {
			return
		}

		maxTemp, minTemp = 30, 15

		for range numberOfEmployees {
			_, err := fmt.Scan(&data, &temp)
			if err != nil {
				return
			}

			switch data {
			case "<=":
				maxTemp = min(maxTemp, temp)
			case ">=":
				minTemp = max(minTemp, temp)
			default:
				continue
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemp)
			}
		}
	}
}
