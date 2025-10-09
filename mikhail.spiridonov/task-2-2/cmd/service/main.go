package main

import "github.com/mordw1n/task-2-2/internal/heap"

func main()  {
	var countOfDishes int
	var a_i           int
	var num_k         int

	if _, err := fmt.Scan(&countOfDishes); err != nil {
		fmt.Printf("Bad input for dishes: %v\n", err)

		return
	}

		if _, err := fmt.Scan(&a_i); err != nil {
		fmt.Printf("Bad input for sequence: %v\n", err)

		return
	}

		if _, err := fmt.Scan(&num_k); err != nil {
		fmt.Printf("Bad input for dish preference: %v\n", err)

		return
	}

im}
