package main

import "fmt"

const (
	MinTemp  = 15
	MaxTemp  = 30
	ErrorVal = -1
)

func main() {
	var departments int
	var employees int
	var temperature int
	var operation string

	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println("Incorrect input")

		return
	}

	for range departments {
		if _, err := fmt.Scan(&employees); err != nil {
			fmt.Println("Incorrect input")

			return
		}

		currentMin := MinTemp
		currentMax := MaxTemp
		validRange := true

		for range employees {
			if _, err := fmt.Scan(&operation, &temperature); err != nil {
				fmt.Println("Incorrect input")

				return
			}

			if !validRange {
				fmt.Println(ErrorVal)

				continue
			}

			switch operation {
			case ">=":
				if temperature > currentMin {
					currentMin = temperature
				}
			case "<=":
				if temperature < currentMax {
					currentMax = temperature
				}
			default:
				fmt.Println("Incorrect input")

				return
			}

			if currentMin > currentMax {
				validRange = false
				fmt.Println(ErrorVal)
			} else {
				fmt.Println(currentMin)
			}
		}
	}
}
