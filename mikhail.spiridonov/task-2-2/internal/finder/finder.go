package finder

import "container/heap"

func finderTheLargest(nums []int, k int) int {
	object := &Heap{}
	heap.Init(object)

	for _, num := range nums {
		heap.Push(object, num)
		if object.Len() > k {
			heap.Pop(object)
		}
	}

	return (*object)[0]
}