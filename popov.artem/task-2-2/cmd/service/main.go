package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(x interface{}) {
	item, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, item)
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	numItems := len(old)
	x := old[numItems-1]
	*h = old[0 : numItems-1]

	return x
}

func main() {
	scanner := bufio.NewReader(os.Stdin)

	nLine, _ := scanner.ReadString('\n')
	nLine = strings.TrimSpace(nLine)
	numDishes, _ := strconv.Atoi(nLine)

	aLine, _ := scanner.ReadString('\n')
	aLine = strings.TrimSpace(aLine)
	parts := strings.Split(aLine, " ")

	nums := make([]int, numDishes)
	for i, part := range parts {
		nums[i], _ = strconv.Atoi(part)
	}

	kLine, _ := scanner.ReadString('\n')
	kLine = strings.TrimSpace(kLine)
	kth, _ := strconv.Atoi(kLine)

	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)

	for _, num := range nums {
		heap.Push(maxHeap, num)
	}

	var result int

	for range kth {
		if value, ok := heap.Pop(maxHeap).(int); ok {
			result = value
		}
	}

	fmt.Println(result)
}
