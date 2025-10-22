package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"container/heap"
)

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	item, ok := x.(int)
	if !ok {
		return
	}
	*h = append(*h, item)
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	scanner := bufio.NewReader(os.Stdin)

	nLine, _ := scanner.ReadString('\n')
	nLine = strings.TrimSpace(nLine)
	n, _ := strconv.Atoi(nLine)

	aLine, _ := scanner.ReadString('\n')
	aLine = strings.TrimSpace(aLine)
	parts := strings.Split(aLine, " ")
	nums := make([]int, n)
	for i, part := range parts {
		nums[i], _ = strconv.Atoi(part)
	}

	kLine, _ := scanner.ReadString('\n')
	kLine = strings.TrimSpace(kLine)
	k, _ := strconv.Atoi(kLine)

	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)

	for _, num := range nums {
		heap.Push(maxHeap, num)
	}

	var result int
	for range k {
		item := heap.Pop(maxHeap)
		result = item.(int)
	}

	fmt.Println(result)
}