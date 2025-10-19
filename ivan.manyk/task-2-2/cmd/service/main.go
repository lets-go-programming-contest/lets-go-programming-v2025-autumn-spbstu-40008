package main

import (
	"container/heap"
	"fmt"
)

type Heap []int

func (h *Heap) Len() int {
	return len(*h)
}

func (h *Heap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *Heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *Heap) Pop() interface{} {
	lenHeap := len(*h)
	x := (*h)[lenHeap-1]
	*h = (*h)[0 : lenHeap-1]

	return x
}

func main() {
	var countOfDishes int
	_, err := fmt.Scan(&countOfDishes)
	if err != nil {
		fmt.Println("Error with count of dishes, code error:", err)

		return
	}

	ratingOfDishes := make([]int, countOfDishes)
	for i := 0; i < countOfDishes; i++ {
		_, err = fmt.Scan(&ratingOfDishes[i])
		if err != nil {
			fmt.Println("Error with rating of dishes, code error:", err)

			return
		}
	}

	var dishNumber int
	_, err = fmt.Scan(&dishNumber)
	if err != nil {
		fmt.Println("Error with prefered dish, code error:", err)

		return
	}

	if dishNumber > countOfDishes {
		fmt.Println("The number of dishes exceeds the count of dishes")

		return
	}

	currentHeap := &Heap{}
	heap.Init(currentHeap)
	for _, rating := range ratingOfDishes {
		if currentHeap.Len() < dishNumber {
			heap.Push(currentHeap, rating)
		} else if rating > (*currentHeap)[0] {
			heap.Pop(currentHeap)
			heap.Push(currentHeap, rating)
		}
	}

	fmt.Println((*currentHeap)[0])
}
