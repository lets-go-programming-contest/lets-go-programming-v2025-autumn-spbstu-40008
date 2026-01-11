package main

import (
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var amountOfDishes, preference, tempRating int

	dishesRatings := &IntHeap{}
	heap.Init(dishesRatings)

	_, err := fmt.Scan(&amountOfDishes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading the number of dishes: %v\n", err)

		return
	}

	if amountOfDishes < 1 {
		fmt.Fprintf(os.Stderr, "the number of dishes must be more than 0\n")
	}

	for range amountOfDishes {
		_, err = fmt.Scan(&tempRating)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading the rating of the dish: %v\n", err)

			return
		}

		heap.Push(dishesRatings, tempRating)
	}

	_, err = fmt.Scan(&preference)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading priority: %v\n", err)

		return
	}

	if preference > amountOfDishes || preference < 1 {
		fmt.Fprintf(os.Stderr, "incorrect preferred dish number\n")
	} else {
		for i := amountOfDishes; i > preference; i-- {
			heap.Pop(dishesRatings)
		}

		fmt.Println(heap.Pop(dishesRatings))
	}
}
