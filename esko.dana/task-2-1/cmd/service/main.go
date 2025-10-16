package main

import (
	"fmt"
)

func main() {
	var (
		numberOfDepartments, numberOfEmployees, minTemp, maxTemp, t int
		operation                                                   string
	)

	_, err := fmt.Scan(&numberOfDepartments)
	if err != nil {
		fmt.Println("Incorrect input")

		return
	}

	for range numberOfDepartments {
		_, err := fmt.Scan(&numberOfEmployees)
		if err != nil {
			fmt.Println("Incorrect input")

			return
		}

		minTemp = 15
		maxTemp = 30

		for range numberOfEmployees {
			_, err := fmt.Scan(&operation, &t)
			if err != nil {
				fmt.Println("Incorrect input")

				return
			}

			switch operation {
			case ">=":
				if t > minTemp {
					minTemp = t
				}
			case "<=":
				if t < maxTemp {
					maxTemp = t
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
