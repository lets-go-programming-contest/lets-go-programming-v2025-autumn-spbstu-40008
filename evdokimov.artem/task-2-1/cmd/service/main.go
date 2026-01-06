package main

import "fmt"

const (
	MinTemp  = 15
	MaxTemp  = 30
	ErrorVal = -1
)

func main() {
	var departments int

	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println("Incorrect input")

		return
	}

	for range departments {
		processDepartment()
	}
}

func processDepartment() {
	var employees int

	if _, err := fmt.Scan(&employees); err != nil {
		fmt.Println("Incorrect input")

		return
	}

	currentMin := MinTemp
	currentMax := MaxTemp
	validRange := true

	for range employees {
		validRange = processEmployee(currentMin, currentMax, validRange)

		if !validRange {
			currentMin = MinTemp
			currentMax = MaxTemp
		}
	}
}

func processEmployee(currentMin, currentMax int, validRange bool) bool {
	var operation string
	var temperature int

	if _, err := fmt.Scan(&operation, &temperature); err != nil {
		fmt.Println("Incorrect input")

		return false
	}

	if !validRange {
		fmt.Println(ErrorVal)

		return false
	}

	newMin, newMax := updateTemperatureRange(operation, temperature, currentMin, currentMax)

	if newMin > newMax {
		fmt.Println(ErrorVal)

		return false
	}

	fmt.Println(newMin)

	return true
}

func updateTemperatureRange(operation string, temperature, min, max int) (int, int) {
	switch operation {
	case ">=":
		if temperature > min {
			min = temperature
		}
	case "<=":
		if temperature < max {
			max = temperature
		}
	default:
		fmt.Println("Incorrect input")
	}

	return min, max
}
