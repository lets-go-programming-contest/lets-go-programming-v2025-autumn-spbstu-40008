package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var numDishes, dishRating, numInTop int

	sortHeap := &IntHeap{}
	heap.Init(sortHeap)

	if _, err := fmt.Scanln(&numDishes); err != nil {
		return
	}

	for range numDishes {
		if _, err := fmt.Scan(&dishRating); err != nil {
			heap.Push(sortHeap, dishRating)
		}
	}

	if _, err := fmt.Scanln(&numInTop); err != nil {
		return
	}

	for sortHeap.Len() > numInTop {
		heap.Pop(sortHeap)
	}

	fmt.Println((*sortHeap)[0])
}
