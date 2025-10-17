package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}
func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}
func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
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
	h := &IntHeap{}
	heap.Init(h)

	for i := 0; i < len(amountOfDishes); i++ {
		if h.Len() < numberk {
			heap.Push(h, amountOfDishes[i])
		} else if amountOfDishes[i] > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, amountOfDishes[i])
		}
	}
	return (*h)[0]
}

func main() {
	var amountOfDishes, numberk int

	fmt.Scan(&amountOfDishes)
	arr := make([]int, amountOfDishes)
	for i := 0; i < amountOfDishes; i++ {
		fmt.Scan(&arr[i])
	}

	fmt.Scan(&numberk)
	fmt.Println(searching(arr, numberk))
}
