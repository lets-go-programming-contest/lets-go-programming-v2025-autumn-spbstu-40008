package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func main() {
	var countFood, currFood, numFood int

	foodHeap := &IntHeap{}
	heap.Init(foodHeap)

	_, err := fmt.Scan(&countFood)
	if err != nil {
		fmt.Println(err)

		return
	}

	for range countFood {
		_, err = fmt.Scan(&currFood)
		if err != nil {
			fmt.Println(err)

			return
		}

		heap.Push(foodHeap, currFood)
	}

	_, err = fmt.Scan(&numFood)
	if err != nil {
		fmt.Println(err)

		return
	}

	for range numFood - 1 {
		heap.Pop(foodHeap)
	}

	fmt.Println(heap.Pop(foodHeap))
}
