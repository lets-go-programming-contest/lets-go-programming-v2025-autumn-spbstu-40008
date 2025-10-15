package myheap

type MyHeap []int

func (object *MyHeap) Len() int {
	if object == nil {
		return 0
	}

	return len(*object)
}

func (object *MyHeap) Less(firstIndex, secondIndex int) bool {
	if object == nil {
		return false
	}

	if firstIndex < 0 || firstIndex >= len(*object) || secondIndex < 0 || secondIndex >= len(*object) {
		return false
	}

	return (*object)[firstIndex] < (*object)[secondIndex]
}

func (object *MyHeap) Swap(firstIndex, secondIndex int) {
	if object == nil {
		return
	}

	if firstIndex < 0 || firstIndex >= len(*object) || secondIndex < 0 || secondIndex >= len(*object) {
		return
	}

	(*object)[firstIndex], (*object)[secondIndex] = (*object)[secondIndex], (*object)[firstIndex]
}

func (object *MyHeap) Push(number interface{}) {
	if object == nil {
		value, ok := number.(int)
		if !ok {
			return
		}

		newHeap := MyHeap{value}
		*object = newHeap

		return
	}

	value, ok := number.(int)
	if !ok {
		return
	}

	*object = append(*object, value)
}

func (object *MyHeap) Pop() interface{} {
	if object == nil || len(*object) == 0 {
		return nil
	}

	outdated := *object
	tempLen := len(outdated)

	if tempLen == 0 {
		return nil
	}

	number := outdated[tempLen-1]
	*object = outdated[0 : tempLen-1]

	return number
}
