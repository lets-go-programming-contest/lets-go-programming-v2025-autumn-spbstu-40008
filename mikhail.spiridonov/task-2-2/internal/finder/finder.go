package finder

import "github.com/mordw1n/task-2-2/internal/heap"

func finderTheLargest(nums []int, k int) int {
	object := &MinHeap{}
	heap.Init(object)

	for _, num := range nums {
		heap.Push(object, num)
		if object.Len() > k {
			heap.Pop(object)
		}
	}

	return (*object)[0]
}