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

	a_i := make([]int, countOfDishes)
	for i := 0; i < countOfDishes; i++ {
		if _, err := fmt.Scan(&a_i[i]); err != nil {
			fmt.Printf("Bad input for sequence: %v\n", err)

			return
		}
	}

	var num_k int
	if _, err := fmt.Scan(&num_k); err != nil {
		fmt.Printf("Bad input for dish preference: %v\n", err)

		return
	}

	result := finder.finderTheLargest(a_i, num_k)
	fmt.Println(result)
}