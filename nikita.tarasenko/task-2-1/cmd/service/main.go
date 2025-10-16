package main

import "fmt"

const (
	minimum_temperature = 15
	maximum_temperature = 30
)

func main() {
	var department_count int
	var employee_count int

	var now_minimum_temperature int
	var now_maximum_temperature int

	var data int
	var sign string

	_, err := fmt.Scan(&department_count)
	if err != nil {
		fmt.Printf("Try another count of department\n")
		return
	}

	for range department_count {
		_, err := fmt.Scan(&employee_count)
		if err != nil {
			fmt.Printf("Try another count of employee\n")
			return
		}

		now_minimum_temperature = minimum_temperature
		now_maximum_temperature = maximum_temperature
		for range employee_count {
			_, err := fmt.Scan(&sign, &data)
			if err != nil {
				fmt.Printf("Bad format for sign and data\n")
				return
			}

			switch sign {
			case ">=":
				if data > now_minimum_temperature {
					now_minimum_temperature = data
				}
			case "<=":
				if data < now_maximum_temperature {
					now_maximum_temperature = data
				}
			default:
				fmt.Printf("Try another sign\n")
				return
			}

			if now_minimum_temperature > now_maximum_temperature {
				fmt.Printf("-1\n")
				continue
			}
			fmt.Printf("%d\n", now_minimum_temperature)
		}
	}
}
