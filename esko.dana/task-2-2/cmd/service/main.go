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

func (h *IntHeap) Push(x any) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var dishNum, preference, temp int

	rating := &IntHeap{}
	heap.Init(rating)

	_, err := fmt.Scan(&dishNum)
	if err != nil {
		fmt.Println("Invalid number")

		return
	}

	for range dishNum {
		_, err = fmt.Scan(&temp)
		if err != nil {
			fmt.Println("Invalid number")

			return
		}

		heap.Push(rating, temp)
	}

	_, err = fmt.Scan(&preference)
	if err != nil {
		fmt.Println("Invalid number")

		return
	}

	for range preference - 1 {
		heap.Pop(rating)
	}

	fmt.Println(heap.Pop(rating))
}
