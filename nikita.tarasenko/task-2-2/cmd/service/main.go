package main

import (
	"container/heap"
	"fmt"
)

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	value, err := x.(int)
	if !err {
		fmt.Println("push error:", err)

		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func main() {
	var dishesCount, dish, preference int

	DishHeap := &IntHeap{}
	heap.Init(DishHeap)

	_, err := fmt.Scan(&dishesCount)
	if err != nil {
		fmt.Println("Try another number of dishes")

		return
	}

	for range dishesCount {
		_, err := fmt.Scan(&dish)
		if err != nil {
			fmt.Println("Try another dish")

			return
		}

		heap.Push(DishHeap, dish)
	}

	_, err = fmt.Scan(&preference)
	if err != nil {
		fmt.Println("Try another preference")

		return
	}

	for range preference - 1 {
		heap.Pop(DishHeap)
	}

	fmt.Println((*DishHeap)[0])
}
