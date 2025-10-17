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
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func (h IntHeap) Length() int {
	return h.Len()
}

func (h IntHeap) Lessthan(i, j int) bool {
	return h.Less(i, j)
}

func searching(amountOfDishes []int, numberk int) int {
	heapExample := &IntHeap{}
	heap.Init(heapExample)

	for _, dish := range amountOfDishes {
		if heapExample.Len() < numberk {
			heap.Push(heapExample, dish)
		} else if dish > (*heapExample)[0] {
			heap.Pop(heapExample)
			heap.Push(heapExample, dish)
		}
	}

	return (*heapExample)[0]
}

func main() {
	var amountOfDishes, numberk int

	_, err := fmt.Scan(&amountOfDishes)
	if err != nil {
		return
	}

	arr := make([]int, amountOfDishes)
	for i := range arr {
		_, err = fmt.Scan(&arr[i])
		if err != nil {
			return
		}
	}

	_, err = fmt.Scan(&numberk)
	if err != nil {
		return
	}

	result := searching(arr, numberk)
	fmt.Println(result)
}
