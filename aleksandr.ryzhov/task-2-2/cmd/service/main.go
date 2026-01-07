package main

import (
	"container/heap"
	"fmt"
)

type intHeap []int

func (h *intHeap) Len() int           { return len(*h) }
func (h *intHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *intHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
func (h *intHeap) Push(x any) {
	num, err := x.(int)
	if !err {
		fmt.Println("Incorrect number of departments", err)

		return
	}
	*h = append(*h, num)
}
func (h *intHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var dishesCount, prefDishNum, tmp int

	// Получаем количество блюд
	_, err := fmt.Scan(&dishesCount)
	if err != nil || dishesCount < 1 {
		fmt.Println("Incorrect number of meals ", err)

		return
	}

	dishes := intHeap{}
	heap.Init(&dishes)

	// Получаем блюда по приоритетам
	for range dishesCount {
		_, err := fmt.Scan(&tmp)
		if err != nil {
			fmt.Println("Incorrect priority ", err)

			return
		}
		heap.Push(&dishes, tmp)
	}

	// Получаем приоритет предпочитаемого блюда
	_, err = fmt.Scan(&prefDishNum)
	if err != nil {
		fmt.Println("Incorrect priority ", err)

		return
	}

	// Выбираем блюдо с нужным приоритетом
	for range dishes {
		if dishes.Len() == prefDishNum {
			break
		}
		heap.Pop(&dishes)
	}
	fmt.Println(heap.Pop(&dishes))
}
