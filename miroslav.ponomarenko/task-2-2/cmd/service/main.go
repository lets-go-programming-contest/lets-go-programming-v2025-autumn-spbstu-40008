package main

import (
	"container/heap"
	"fmt"
)

type Heap []int

func (h *Heap) Len() int           { return len(*h) }
func (h *Heap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *Heap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *Heap) Push(x any) {
	if v, ok := x.(int); ok {
		*h = append(*h, v)
	}
}

func (h *Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var dishes, dish, prefer int

	myHeap := &Heap{}
	heap.Init(myHeap)

	if _, err := fmt.Scanln(&dishes); err != nil {
		fmt.Println("Invalid dishes value")

		return
	}

	for range dishes {
		if _, err := fmt.Scan(&dish); err != nil {
			fmt.Println("Invalid dish value")

			return
		}

		heap.Push(myHeap, dish)
	}

	if _, err := fmt.Scanln(&prefer); err != nil {
		fmt.Println("Invalid prefered dish value")

		return
	}

	for myHeap.Len() > prefer {
		heap.Pop(myHeap)
	}

	fmt.Println((*myHeap)[0])
}
