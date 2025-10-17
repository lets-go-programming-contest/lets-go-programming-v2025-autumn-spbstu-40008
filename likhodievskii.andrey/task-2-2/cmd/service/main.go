package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var countDishes, numberWishDish int

	if _, err := fmt.Scan(&countDishes); err != nil {
		fmt.Printf("Bad count dishes value: %v\n", err)

		return

	}

	if countDishes < 0 {
		fmt.Printf("Bad count dishes value: %d\n", countDishes)

		return
	}

	dishes := make([]int, countDishes)

	for indexDishesNow := 0; indexDishesNow < countDishes; indexDishesNow++ {
		if _, err := fmt.Scan(&dishes[indexDishesNow]); err != nil {
			fmt.Printf("Bad point for dish: %v\n", err)

			return
		}
	}

	if _, err := fmt.Scan(&numberWishDish); err != nil {
		fmt.Printf("Bad value for number dish for customer: %v\n", err)

		return
	}

	if numberWishDish < 1 || numberWishDish >= countDishes {
		fmt.Printf("Unreacheble number dish for customer: %d\n", numberWishDish)

		return
	}

	heapForDishes := &IntHeap{}
	heap.Init(heapForDishes)

	for indexFirstNumbersWishDish := 0; indexFirstNumbersWishDish < numberWishDish; indexFirstNumbersWishDish++ {
		heap.Push(heapForDishes, dishes[indexFirstNumbersWishDish])
	}

	for countDishBeforeNeeded := numberWishDish; countDishBeforeNeeded < countDishes; countDishBeforeNeeded++ {
		if dishes[countDishBeforeNeeded] > (*heapForDishes)[0] {
			heap.Pop(heapForDishes)
			heap.Push(heapForDishes, dishes[countDishBeforeNeeded])
		}
	}

	fmt.Println((*heapForDishes)[0])
}
