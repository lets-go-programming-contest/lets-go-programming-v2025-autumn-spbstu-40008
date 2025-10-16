package main

import "fmt"

func main() {
	var numDepartments, numEmployees int

	if _, err := fmt.Scan(&numDepartments); err != nil {
		return
	}

	for range numDepartments {
		if _, err := fmt.Scan(&numEmployees); err != nil {
			return
		}

		minTemp, maxTemp := 15, 30

		for range numEmployees {
			var operator string
			var temp int

			if _, err := fmt.Scan(&operator, &temp); err != nil {
				return
			}

			switch operator {
			case ">=":
				if temp > minTemp {
					minTemp = temp
				}
			case "<=":
				if temp < maxTemp {
					maxTemp = temp
				}
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
