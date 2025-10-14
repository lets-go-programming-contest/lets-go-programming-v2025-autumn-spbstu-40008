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
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x any) {
	temp, ok := x.(int)
	if !ok {
		panic("Expected int in heap")
	}

	*h = append(*h, temp)
}

func (h *IntHeap) Pop() any {
	n := len(*h)
	x := (*h)[n-1]
	*h = (*h)[0 : n-1]

	return x
}

func main() {
	var dishesNumber int

	if _, err := fmt.Scan(&dishesNumber); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read number of dishes: %v\n", err)

		return
	}

	var dishesHeap IntHeap

	for range dishesNumber {
		var currentDishValue int

		if _, err := fmt.Scan(&currentDishValue); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to read dish value: %v\n", err)

			return
		}

		heap.Push(&dishesHeap, currentDishValue)
	}

	var preferredDishNumber int

	if _, err := fmt.Scan(&preferredDishNumber); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read preferred dish number: %v\n", err)

		return
	}

	if preferredDishNumber <= 0 || preferredDishNumber > dishesHeap.Len() {
		fmt.Fprintf(os.Stderr, "Error: invalid preferred dish number\n")

		return
	}

	for range preferredDishNumber - 1 {
		heap.Pop(&dishesHeap)
	}

	fmt.Println(heap.Pop(&dishesHeap))
}
