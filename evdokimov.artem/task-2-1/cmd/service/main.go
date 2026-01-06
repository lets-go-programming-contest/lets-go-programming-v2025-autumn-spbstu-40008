package main

import "fmt"

const (
	MinTemp  = 15
	MaxTemp  = 30
	ErrorVal = -1
)

func main() {
	var departments, employees, temperature int
	var operation string

	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println("Incorrect input")
		return
	}

	for i := 0; i < departments; i++ {
		if _, err := fmt.Scan(&employees); err != nil {
			fmt.Println("Incorrect input")
			return
		}

		currentMin := MinTemp
		currentMax := MaxTemp
		validRange := true

		for j := 0; j < employees; j++ {
			if _, err := fmt.Scan(&operation, &temperature); err != nil {
				fmt.Println("Incorrect input")
				return
			}

			if validRange {
				if operation == ">=" {
					if temperature > currentMin {
						currentMin = temperature
					}
				} else if operation == "<=" {
					if temperature < currentMax {
						currentMax = temperature
					}
				} else {
					fmt.Println("Incorrect input")
					return
				}

				if currentMin > currentMax {
					validRange = false
				}
			}

			if validRange {
				fmt.Println(currentMin)
			} else {
				fmt.Println(ErrorVal)
			}
		}
	}
}
