package main

import (
	"container/heap"
	"fmt"
	"os"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(elem interface{}) {
	if val, ok := elem.(int); ok {
		*h = append(*h, val)
	}
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	elem := old[n-1]
	*h = old[0 : n-1]

	return elem
}

func main() {
	var countDishes, priority, result int

	dishes := &IntHeap{}
	heap.Init(dishes)

	if _, err := fmt.Scan(&countDishes); err != nil {
		fmt.Println("Error input for count dishes")
		os.Exit(0)
	}

	for range countDishes {
		var valDish int

		if _, err := fmt.Scan(&valDish); err != nil {
			fmt.Println("Error input for value dish")
			os.Exit(0)
		}

		heap.Push(dishes, valDish)
	}

	if _, err := fmt.Scan(&priority); err != nil {
		fmt.Println("Error input for dish's priority")
		os.Exit(0)
	}

	priority = dishes.Len() - priority + 1

	for range priority {
		if val, ok := heap.Pop(dishes).(int); ok {
			result = val
		}
	}

	fmt.Println(result)
}
