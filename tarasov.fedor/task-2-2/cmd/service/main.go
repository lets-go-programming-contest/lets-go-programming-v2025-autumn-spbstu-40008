package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(x interface{}) {
	if val, ok := x.(int); ok {
		*h = append(*h, val)
	}
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var numberOfDishes int

	_, err := fmt.Scanln(&numberOfDishes)
	if err != nil {
		fmt.Println("Input error: Please enter a valid number of dishes.")

		return
	}

	nums := make([]int, numberOfDishes)
	for index := range numberOfDishes {
		if _, err := fmt.Scan(&nums[index]); err != nil {
			fmt.Println("Input error: Please enter a valid sequence of dishes.")

			return
		}
	}

	heapDishes := &MaxHeap{}

	for _, num := range nums {
		heap.Push(heapDishes, num)
	}

	var preferredDishNum int

	_, err = fmt.Scanln(&preferredDishNum)
	if err != nil || preferredDishNum > numberOfDishes || preferredDishNum <= 0 {
		fmt.Println("Invalid number of prefer dish")

		return
	}

	var preferDish interface{}
	for range preferredDishNum {
		preferDish = heap.Pop(heapDishes)
	}

	fmt.Println(preferDish)
}
