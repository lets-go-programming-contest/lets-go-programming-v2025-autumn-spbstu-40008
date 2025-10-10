package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (heap *IntHeap) Len() int {
	return len(*heap)
}

func (heap *IntHeap) Less(i, j int) bool {
	return (*heap)[i] < (*heap)[j]
}

func (heap *IntHeap) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *IntHeap) Push(value any) {
	intValue, ok := value.(int)
	if !ok {
		panic("There is no int in the heap")
	}

	*heap = append(*heap, intValue)
}

func (heap *IntHeap) Pop() any {
	length := len(*heap)
	lastElement := (*heap)[length-1]
	*heap = (*heap)[0 : length-1]

	return lastElement
}

func main() {
	var (
		countDish    int
		priorityDish int
	)

	_, err := fmt.Scan(&countDish)
	if err != nil || countDish < 1 {
		fmt.Println(-1)

		return
	}

	Rate := make([]int, countDish)
	for i := range countDish {
		_, err := fmt.Scan(&Rate[i])
		if err != nil {
			fmt.Println(-1)

			return
		}
	}

	_, err = fmt.Scan(&priorityDish)
	if err != nil || priorityDish < 1 || priorityDish > countDish {
		fmt.Println(-1)

		return
	}

	result := getPriority(Rate, priorityDish)
	fmt.Println(result)
}

func getPriority(dishes []int, priority int) int {
	heapData := &IntHeap{}
	heap.Init(heapData)

	for _, rating := range dishes {
		heap.Push(heapData, rating)

		if heapData.Len() > priority {
			heap.Pop(heapData)
		}
	}

	return (*heapData)[0]
}
