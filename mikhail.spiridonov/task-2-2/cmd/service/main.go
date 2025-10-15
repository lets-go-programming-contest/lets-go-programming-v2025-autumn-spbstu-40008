package main

import (
 	"fmt"

 	"github.com/mordw1n/task-2-2/internal/finder"
)

func main()  {
	var countOfDishes int
	if _, err := fmt.Scan(&countOfDishes); err != nil {
		fmt.Printf("Bad input for dishes: %v\n", err)

		return
	}

	aIth := make([]int, countOfDishes)
	for i := range countOfDishes {
		if _, err := fmt.Scan(&aIth[i]); err != nil {
			fmt.Printf("Bad input for sequence: %v\n", err)

			return
		}
	}

	var numKth int
	if _, err := fmt.Scan(&numKth); err != nil {
		fmt.Printf("Bad input for dish preference: %v\n", err)

		return
	}

	result := finder.FinderTheLargest(aIth, numKth)
	fmt.Println(result)
}