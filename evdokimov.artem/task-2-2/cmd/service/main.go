package main

import (
	"container/heap"
	"fmt"
)

type RatingHeap []int

func (h RatingHeap) Len() int {
	return len(h)
}

func (h RatingHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h RatingHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *RatingHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *RatingHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]

	return item
}

func main() {
	var dishCount, k int

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

	_, err = fmt.Scan(&k)
	if err != nil || k < 1 || k > dishCount {
		fmt.Println(-1)

		return
	}

	for i := 1; i < k; i++ {
		heap.Pop(ratings)
	}

	result := heap.Pop(ratings).(int)
	fmt.Println(result)
}
