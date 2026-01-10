package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type DishesHeap []int

func (h DishesHeap) Len() int { return len(h) }

func (heapSlice DishesHeap) Less(i, j int) bool {
	return heapSlice[i] > heapSlice[j]
}

func (heapSlice DishesHeap) Swap(i, j int) {
	heapSlice[i], heapSlice[j] = heapSlice[j], heapSlice[i]
}

func (heapSlice *DishesHeap) Push(x any) {
	*heapSlice = append(*heapSlice, x.(int))
}

func (heapSlice *DishesHeap) Pop() any {
	oldSlice := *heapSlice
	length := len(oldSlice)
	element := oldSlice[length-1]
	*heapSlice = oldSlice[0 : length-1]
	return element
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	defer func() {
		_ = writer.Flush()
	}()

	var totalDishes int
	if _, err := fmt.Fscan(reader, &totalDishes); err != nil {
		return
	}

	dishes := make(DishesHeap, 0, totalDishes)

	for range totalDishes {
		var preferenceValue int
		if _, err := fmt.Fscan(reader, &preferenceValue); err != nil {
			break
		}
		dishes = append(dishes, preferenceValue)
	}

	heap.Init(&dishes)

	var targetRank int
	if _, err := fmt.Fscan(reader, &targetRank); err != nil {
		return
	}

	var result int

	for range targetRank {
		if poppedElement, ok := heap.Pop(&dishes).(int); ok {
			result = poppedElement
		}
	}

	_, _ = fmt.Fprintln(writer, result)
}
