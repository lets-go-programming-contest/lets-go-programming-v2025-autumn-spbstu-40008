package myheap

type myHeap []int

func (object myHeap) Len() int {
	return len(object)
} 

func (object myHeap) Less(i, j int) bool {
	return object[i] < object[j]
}

func (object myHeap) Swap(i, j int) {
	object[i], object[j] = object[j], object[i]
}

func (object *myHeap) Push(number interface{}) {
	*object = append(*object, number.(int))
}

func (object *myHeap) Pop() interface{} {
	outdated := *object
	temp_len := len(outdated)
	number := outdated[temp_len - 1]
	*object = outdated[0 : temp_len - 1]
	
	return number
}