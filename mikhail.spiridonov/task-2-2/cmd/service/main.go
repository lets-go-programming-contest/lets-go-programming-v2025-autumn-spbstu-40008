package main

import (
	"fmt"

	"github.com/mordw1n/task-2-2/internal/finder"
)

func main() {
	var countOfDishes int

	if _, err := fmt.Scan(&countOfDishes); err != nil {
		fmt.Printf("Bad input for dishes: %v\n", err)

		return
	}

	if countOfDishes < 1 || countOfDishes > 10000 {
		fmt.Printf("Invalid number of dishes: %d. Must be between 1 and 10000.\n", countOfDishes)

		return
	}

	aIth := make([]int, countOfDishes)
	for index := range countOfDishes {
		if _, err := fmt.Scan(&aIth[index]); err != nil {
			fmt.Printf("Bad input for sequence: %v\n", err)

			return
		}

		if aIth[index] < -10000 || aIth[index] > 10000 {
			fmt.Printf("Invalid a_i value: %d. Must be between -10000 and 10000.\n", aIth[index])

			return
		}
	}

	var numKth int
	if _, err := fmt.Scan(&numKth); err != nil {
		fmt.Printf("Bad input for dish preference: %v\n", err)

		return
	}

	if numKth < 1 || numKth > countOfDishes {
		fmt.Printf("Invalid k value: %d. Must be between 1 and %d.\n", numKth, countOfDishes)

		return
	}

	result := finder.FinderTheLargest(aIth, numKth)
	fmt.Println(result)
}
