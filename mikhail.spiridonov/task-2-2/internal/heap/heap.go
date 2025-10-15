package heap

type Heap []int

func (object Heap) Len() int {
	return len(object)
} 

func (object Heap) Less(i, j int) bool {
	return object[i] < object[j]
}

func (object Heap) Swap(i, j int) {
	object[i], object[j] = object[j], object[i]
}

func (object *Heap) Push(number interface{}) {
	*object = append(*object, number.(int))
}

func (object *Heap) Pop() interface{} {
	outdated := *object
	temp_len := len(outdated)
	number := outdated[temp_len - 1]
	*object = outdated[0 : temp_len - 1]

	return number
}