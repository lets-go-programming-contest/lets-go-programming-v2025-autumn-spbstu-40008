package main

import (
	"fmt"
	"os"
)

func main() {
	var departments, employees, minT, maxT uint16

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading the number of departments: %v\n", err)

		return
	}

	for range departments {
		maxT = 30
		minT = 15

		_, err = fmt.Scan(&employees)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error readind the number of employees: %v\n", err)

			return
		}

		for range employees {
			var operator string

			var temperature uint16

			_, err = fmt.Scan(&operator)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading the logical operator: %v", err)

				return
			}

			_, err = fmt.Scan(&temperature)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading temperature: %v", err)

				return
			}

			switch operator {
			case ">=":
				minT = max(minT, temperature)
			case "<=":
				maxT = min(maxT, temperature)
			default:
				fmt.Fprintf(os.Stderr, "use only the operators >= or <=")
			}

			if minT > maxT {
				fmt.Println(-1)
			} else {
				fmt.Println(minT)
			}
		}
	}
}
