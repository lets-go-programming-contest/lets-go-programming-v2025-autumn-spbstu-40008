package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"container/heap"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
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

	h := &IntHeap{}
	heap.Init(h)

	for _, num := range nums {
		heap.Push(h, num)
	}

	var result int
	for i := 0; i < k; i++ {
		result = heap.Pop(h).(int)
	}

	fmt.Println(result)
}