package main

import "fmt"

const (
	minimumTemperature = 15
	maximumTemperature = 30
)

func main() {
	var departmentCount int

	var employeeCount int

	var nowMinimumTemperature int

	var nowMaximumTemperature int

	var data int

	var sign string

	_, err := fmt.Scan(&departmentCount)
	if err != nil {
		fmt.Printf("Try another count of department\n")

		return
	}

	for range departmentCount {
		_, err := fmt.Scan(&employeeCount)
		if err != nil {
			fmt.Printf("Try another count of employee\n")

			return
		}

		nowMinimumTemperature = minimumTemperature
		nowMaximumTemperature = maximumTemperature

		for range employeeCount {
			_, err := fmt.Scan(&sign, &data)
			if err != nil {
				fmt.Printf("Bad format for sign and data\n")

				return
			}

			switch sign {
			case ">=":
				if data > nowMinimumTemperature {
					nowMinimumTemperature = data
				}
			case "<=":
				if data < nowMaximumTemperature {
					nowMaximumTemperature = data
				}
			default:
				fmt.Printf("Try another sign\n")

				return
			}

			if nowMinimumTemperature > nowMaximumTemperature {
				fmt.Printf("-1\n")

				continue
			}

			fmt.Printf("%d\n", nowMinimumTemperature)
		}
	}
}
