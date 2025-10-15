package main

import (
	"container/heap"
 	"fmt"
 	"github.com/mordw1n/task-2-2/internal/heap"
 	"github.com/mordw1n/task-2-2/internal/finder"
)

func main()  {
	var countOfDishes int
	if _, err := fmt.Scan(&countOfDishes); err != nil {
		fmt.Printf("Bad input for dishes: %v\n", err)

		return
	}

		var a_i int
		if _, err := fmt.Scan(&a_i); err != nil {
		fmt.Printf("Bad input for sequence: %v\n", err)

		return
	}
		var num_k int
		if _, err := fmt.Scan(&num_k); err != nil {
		fmt.Printf("Bad input for dish preference: %v\n", err)

		return
	}

	result := finderTheLargest

}
