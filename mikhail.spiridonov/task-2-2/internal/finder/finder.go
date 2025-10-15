package finder

import (
	"container/heap"

 	"github.com/mordw1n/task-2-2/internal/myheap"
)

func finderTheLargest(nums []int, k int) int {
	object := &myheap.myHeap{}
	heap.Init(object)

	for _, num := range nums {
		heap.Push(object, num)
		if object.Len() > k {
			heap.Pop(object)
		}
	}

	return (*object)[0]
}