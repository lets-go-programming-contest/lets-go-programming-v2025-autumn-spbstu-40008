package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int            { return len(*h) }
func (h *IntHeap) Less(i, j int) bool  { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)       { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
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
	var count, kth int
	if _, err := fmt.Scan(&count); err != nil {
		return
	}
	arr := make([]int, count)
	for i := range arr {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			return
		}
	}
	if _, err := fmt.Scan(&kth); err != nil {
		return
	}
	minHeap := &IntHeap{}
	heap.Init(minHeap)
	for _, val := range arr {
		heap.Push(minHeap, val)

		if minHeap.Len() > kth {
			heap.Pop(minHeap)
		}
	}

	top := (*minHeap)[0]

	fmt.Println(top)
}
