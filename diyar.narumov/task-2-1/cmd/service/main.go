package main

import "fmt"

func processEmployee(low, high *int) {
	var cmp string

	if _, err := fmt.Scan(&cmp); err != nil {
		fmt.Println("Error: failed to read comparison operator")

		return
	}

	var temperature int

	if _, err := fmt.Scan(&temperature); err != nil {
		fmt.Println("Error: failed to read temperature value")

		return
	}

	switch cmp {
	case "<=":
		if *high > temperature {
			*high = temperature
		}
	case ">=":
		if *low < temperature {
			*low = temperature
		}
	default:
		fmt.Println("Error: invalid comparison operator. Use '<=' or '>='")

		return
	}

	if *low > *high {
		fmt.Println(-1)
	} else {
		fmt.Println(*low)
	}
}

func processDepartment() {
	var employeesNumber int

	if _, err := fmt.Scan(&employeesNumber); err != nil {
		fmt.Println("Error: failed to read number of employees")

		return
	}

	low, high := 15, 30
	for range employeesNumber {
		processEmployee(&low, &high)
	}
}

func main() {
	var departmentsNumber int

	if _, err := fmt.Scan(&departmentsNumber); err != nil {
		fmt.Println("Error: failed to read number of departments")

		return
	}

	for range departmentsNumber {
		processDepartment()
	}
}
