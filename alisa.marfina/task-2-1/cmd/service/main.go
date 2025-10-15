package main

import "fmt"

func main() {
	var (
		numberOfDepartments, numberOfEmployees, maxTemp, minTemp, temp int
		data                                                           string
	)

	_, err := fmt.Scan(&numberOfDepartments)
	if err != nil {
		fmt.Println("Error with numbers of departments", err)
		return
	}

	for range numberOfDepartments {
		_, err := fmt.Scan(&numberOfEmployees)
		if err != nil {
			fmt.Println("Error with numbers of employees", err)
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
				fmt.Println("Error with operator")
				continue
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
			}
			fmt.Println(minTemp)
		}
	}
}
