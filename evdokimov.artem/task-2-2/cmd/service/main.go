package main

import (
	"container/heap"
	"fmt"
)

type RatingHeap []int

func (h *RatingHeap) Len() int {
	return len(*h)
}

func (h *RatingHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *RatingHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *RatingHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		panic("invalid type for RatingHeap")
	}

	*h = append(*h, value)
}

func (h *RatingHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]

	return item
}

func main() {
	var dishCount, position int

	_, err := fmt.Scan(&dishCount)
	if err != nil || dishCount < 1 {
		fmt.Println(-1)

		return
	}

	ratings := &RatingHeap{}
	heap.Init(ratings)

	var rating int
	for range dishCount {
		_, err = fmt.Scan(&rating)
		if err != nil {
			fmt.Println(-1)

			return
		}

		heap.Push(ratings, rating)
	}

	_, err = fmt.Scan(&position)
	if err != nil || position < 1 || position > dishCount {
		fmt.Println(-1)

		return
	}

	for i := 1; i < position; i++ {
		heap.Pop(ratings)
	}

	result, ok := heap.Pop(ratings).(int)
	if !ok {
		fmt.Println(-1)
		return
	}

	fmt.Println(result)
}
