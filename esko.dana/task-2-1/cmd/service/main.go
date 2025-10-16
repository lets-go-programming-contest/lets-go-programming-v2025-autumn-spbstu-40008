package main

import (
	"fmt"
)

func main() {
	var numberOfDepartments int

	_, err := fmt.Scan(&numberOfDepartments)
	if err != nil {

		fmt.Println("Incorrect input")
		return
	}

	for range numberOfDepartments {
		var numberOfEmployees int

		_, err = fmt.Scan(&numberOfEmployees)
		if err != nil {
			fmt.Println("Incorrect input")

			return
		}

		minTemp := 15
		maxTemp := 30

		for range numberOfEmployees {
			var op string
			var t int

			_, err := fmt.Scan(&op, &t)
			if err != nil {
				fmt.Println("Incorrect input")
				return
			}

			switch op {
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
				return
			}

			if minTemp > maxTemp {
				fmt.Println(-1)

				continue
			}

			fmt.Println(minTemp)

		}
	}
}
