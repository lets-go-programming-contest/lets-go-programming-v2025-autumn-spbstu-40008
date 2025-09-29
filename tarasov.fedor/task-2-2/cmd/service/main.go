package main

import (
	"container/heap"
	"fmt"
	"os"
)

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var n int
	_, err := fmt.Scanln(&n)
	if err != nil {
		fmt.Println("Invalid number of dishes")
		os.Exit(0)
	}

	nums := make([]int, n)
	for index := range n {
		if _, err := fmt.Scan(&nums[index]); err != nil {

			return
		}
	}

	ai := &MaxHeap{}

	for _, num := range nums {
		heap.Push(ai, num)
	}

	var k int
	_, err = fmt.Scanln(&k)
	if err != nil || k > n || k <= 0 {
		fmt.Println("Invalid number of prefer dish")
		os.Exit(0)
	}

	var preferDish int
	for i := 0; i < k; i++ {
		preferDish = heap.Pop(ai).(int)
	}

	fmt.Println(preferDish)

}
