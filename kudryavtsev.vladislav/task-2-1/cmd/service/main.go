package main

import (
	"fmt"
)

func processConditions(numConditions int) {
	lowerBound := 15
	upperBound := 30

	for i := 0; i < numConditions; i++ {
		var (
			operator string
			value    int
		)

		if _, err := fmt.Scan(&operator, &value); err != nil {
			return
		}

		switch operator {
		case "<=":
			if lowerBound != -1 {
				if upperBound >= value {
					upperBound = value
				}

				if upperBound < lowerBound {
					lowerBound = -1
				}
			}

		case ">=":
			if lowerBound != -1 {
				if lowerBound <= value {
					lowerBound = value
				}

				if lowerBound > upperBound {
					lowerBound = -1
				}
			}
		}

		fmt.Println(lowerBound)
	}
}

func main() {
	var numTests int

	if _, err := fmt.Scan(&numTests); err != nil {
		return
	}

	for i := 0; i < numTests; i++ {
		var numConditions int

		if _, err := fmt.Scan(&numConditions); err != nil {
			return
		}

		processConditions(numConditions)
	}
}
