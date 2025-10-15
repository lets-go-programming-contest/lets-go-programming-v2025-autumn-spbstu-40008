package finder

import (
	"container/heap"

	"github.com/mordw1n/task-2-2/internal/myheap"
)

func FinderTheLargest(nums []int, kthOrder int) int {
	object := &myheap.MyHeap{}
	heap.Init(object)

	for _, num := range nums {
		heap.Push(object, num)
		if object.Len() > kthOrder {
			heap.Pop(object)
		}
	}

	return (*object)[0]
}
