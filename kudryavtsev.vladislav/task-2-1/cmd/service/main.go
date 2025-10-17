package main

import (
	"fmt"
)

func processConditions(numConditions int) {
	var (
		lowerBound = 15
		upperBound = 30
	)

	for range numConditions {
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

	for range numTests {
		var numConditions int
		if _, err := fmt.Scan(&numConditions); err != nil {
			return
		}

		processConditions(numConditions)
	}
}
