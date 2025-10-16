package main

import (
	"fmt"
)

func main() {
	var (
		numberOfDepartments, numberOfEmployees, minTemp, maxTemp, temperature int
		operation                                                             string
	)

	_, err := fmt.Scan(&numberOfDepartments)
	if err != nil {
		fmt.Println("Incorrect input")

		return
	}

	for range numberOfDepartments {
		_, err = fmt.Scan(&numberOfEmployees)
		if err != nil {
			fmt.Println("Incorrect input")

			return
		}

		minTemp = 15
		maxTemp = 30

		for range numberOfEmployees {
			_, err = fmt.Scan(&operation, &temperature)
			if err != nil {
				fmt.Println("Incorrect input")

				return
			}

			switch operation {
			case ">=":
				if temperature > minTemp {
					minTemp = temperature
				}
			case "<=":
				if temperature < maxTemp {
					maxTemp = temperature
				}
			default:
				fmt.Println("Incorrect input")
			}

			if minTemp > maxTemp {
				fmt.Println(-1)

				continue
			}

			fmt.Println(minTemp)

		}
	}
}
