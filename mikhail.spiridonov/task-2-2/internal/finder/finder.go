package finder

import (
	"container/heap"

	"github.com/mordw1n/task-2-2/internal/myheap"
)

func FinderTheLargest(nums []int, kthOrder int) int {
	if len(nums) == 0 {
		return 0
	}

	if kthOrder <= 0 || kthOrder > len(nums) {
		return 0
	}

	object := &myheap.MyHeap{}
	heap.Init(object)

	for _, num := range nums {
		heap.Push(object, num)

		if object.Len() > kthOrder {
			heap.Pop(object)
		}
	}

	if object.Len() == 0 {
		return 0
	}

	return (*object)[0]
}
