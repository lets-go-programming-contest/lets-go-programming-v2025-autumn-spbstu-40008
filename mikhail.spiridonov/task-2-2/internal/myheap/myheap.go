package myheap

type MyHeap []int

func (object MyHeap) Len() int {
	return len(object)
} 

func (object MyHeap) Less(i, j int) bool {
	return object[i] < object[j]
}

func (object MyHeap) Swap(i, j int) {
	object[i], object[j] = object[j], object[i]
}

func (object *MyHeap) Push(number interface{}) {
	*object = append(*object, number.(int))
}

func (object *MyHeap) Pop() interface{} {
	outdated := *object
	temp_len := len(outdated)
	number := outdated[temp_len - 1]
	*object = outdated[0 : temp_len - 1]
	
	return number
}