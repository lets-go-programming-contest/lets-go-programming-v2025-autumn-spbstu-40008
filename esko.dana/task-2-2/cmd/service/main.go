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
	var dishNum, k, temp int

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

	_, err = fmt.Scan(&k)
	if err != nil {
		fmt.Println("Invalid number")

		return
	}

	for i := 0; i < k-1; i++ {
		heap.Pop(rating)
	}

	fmt.Println(heap.Pop(rating))
}
