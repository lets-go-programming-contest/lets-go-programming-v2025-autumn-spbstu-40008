package main

import (
	"fmt"
)

func processDepartment(deptNum, staffCount int) {
	maxtemp := 30
	mintemp := 15

	for employeeIndex := 1; employeeIndex <= staffCount; employeeIndex++ {
		fmt.Printf("Enter operator and temperature (<= or >= value) employee %d department %d:\n", employeeIndex, deptNum)

		var temperatureData string

		var degrees int

		if _, err := fmt.Scan(&temperatureData, &degrees); err != nil {
			panic(err)
		}

		if degrees < 15 || degrees > 30 {
			panic("Temperature out of allowed range")
		}

		if temperatureData != "<=" && temperatureData != ">=" {
			panic("Invalid operator")
		}

		if temperatureData == "<=" && degrees < maxtemp {
			maxtemp = degrees
		} else if temperatureData == ">=" && degrees > mintemp {
			mintemp = degrees
		}

		if mintemp > maxtemp {
			fmt.Printf("Department %d after employee %d: -1\n", deptNum, employeeIndex)
		} else {
			fmt.Printf("Department %d after employee %d: %d\n", deptNum, employeeIndex, mintemp)
		}
	}
}

func main() {
	fmt.Println("Enter number of departments:")

	var departments int

	if _, err := fmt.Scan(&departments); err != nil {
		panic(err)
	}

	if departments < 1 || departments > 1000 {
		panic("Departments count out of range")
	}

	fmt.Println("Enter number of employees:")

	var staffCount int

	if _, err := fmt.Scan(&staffCount); err != nil {
		panic(err)
	}

	if staffCount < 1 || staffCount > 1000 {
		panic("Employees count out of range")
	}

	for deptNum := 1; deptNum <= departments; deptNum++ {
		processDepartment(deptNum, staffCount)
	}
}
