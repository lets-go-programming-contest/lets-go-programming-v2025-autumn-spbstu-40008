package main

import "fmt"

func main() {
	var numberOfDepartments, numberOfPeople uint

	_, err := fmt.Scan(&numberOfDepartments)
	if err != nil {
		fmt.Println("Error with number of departments, code error: ", err)

		return
	}

	var (
		minTemp, maxTemp uint8
		temp             uint8
		operator         string
	)
	for range numberOfDepartments {
		_, err = fmt.Scan(&numberOfPeople)
		if err != nil {
			fmt.Println("Error with number of people in department, code error: ", err)

			return
		}

		minTemp, maxTemp = 15, 30
		for range numberOfPeople {
			_, err = fmt.Scan(&operator, &temp)
			if err != nil {
				fmt.Println("Error with number of people in department, code error: ", err)

				return
			}

			switch operator {
			case ">=":
				minTemp = max(minTemp, temp)
			case "<=":
				maxTemp = min(maxTemp, temp)
			default:
				fmt.Println("Error with operator")

				return
			}

			if maxTemp < minTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemp)
			}
		}
	}
}
