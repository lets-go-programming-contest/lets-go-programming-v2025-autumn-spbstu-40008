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
	return (*h)[i] < (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	v, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, v)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]

	*h = old[:n-1]

	return x
}

func main() {
	var count int
	if _, err := fmt.Scan(&count); err != nil {
		return
	}

	items := make([]int, count)

	for i := range items {
		if _, err := fmt.Scan(&items[i]); err != nil {
			return
		}
	}

	var kth int
	if _, err := fmt.Scan(&kth); err != nil {
		return
	}

	minHeap := &IntHeap{}

	heap.Init(minHeap)

	for _, val := range items {
		heap.Push(minHeap, val)

		if minHeap.Len() > kth {
			heap.Pop(minHeap)
		}
	}

	if minHeap.Len() == 0 {
		return
	}

	fmt.Println((*minHeap)[0])
}
