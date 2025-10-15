package myheap

type MyHeap []int

func (object *MyHeap) Len() int {
	if object == nil {
		return 0
	}
	return len(*object)
}

func (object *MyHeap) Less(i, j int) bool {
	if object == nil {
		return false
	}
	if i < 0 || i >= len(*object) || j < 0 || j >= len(*object) {
		return false
	}
	return (*object)[i] < (*object)[j]
}

func (object *MyHeap) Swap(i, j int) {
	if object == nil {
		return
	}
	if i < 0 || i >= len(*object) || j < 0 || j >= len(*object) {
		return
	}
	(*object)[i], (*object)[j] = (*object)[j], (*object)[i]
}

func (object *MyHeap) Push(number interface{}) {
	if object == nil {
		newHeap := MyHeap{number.(int)}
		*object = newHeap
		return
	}
	*object = append(*object, number.(int))
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
