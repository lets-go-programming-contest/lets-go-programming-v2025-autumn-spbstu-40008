package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrReadInput       = errors.New("failed to read input")
	ErrPreferenceRange = errors.New("preference out of range")
	ErrHeapOperation   = errors.New("heap operation failed")
)

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	if i < 0 || i >= len(*h) || j < 0 || j >= len(*h) {
		panic("heap index out of range")
	}

	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	if i < 0 || i >= len(*h) || j < 0 || j >= len(*h) {
		panic("heap index out of range")
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("invalid type pushed to IntHeap")
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	length := len(old)

	if length == 0 {
		panic("pop from empty heap")
	}

	val := old[length-1]
	*h = old[:length-1]

	return val
}

func main() {
	var dishCount int

	if _, err := fmt.Scan(&dishCount); err != nil {
		fmt.Println("Failed to read dish count:", err)

		return
	}

	dishRatings := &IntHeap{}
	heap.Init(dishRatings)

	for range dishCount {
		var rating int
		if _, err := fmt.Scan(&rating); err != nil {
			fmt.Println("Failed to read dish rating:", err)

			return
		}

		heap.Push(dishRatings, rating)
	}

	var dishPreference int

	if _, err := fmt.Scan(&dishPreference); err != nil {
		fmt.Println("Failed to read dish preference:", err)

		return
	}

	if dishPreference <= 0 || dishPreference > dishRatings.Len() {
		fmt.Println("Preference out of range")

		return
	}

	for range dishPreference - 1 {
		heap.Pop(dishRatings)
	}

	result := heap.Pop(dishRatings)

	fmt.Println(result)
}
