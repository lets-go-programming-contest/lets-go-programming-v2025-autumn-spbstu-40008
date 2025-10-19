package main

import (
	"container/heap"
	"errors"
	"fmt"

	"timur.akhmetkhanov/task-2-2/internal/intheap"
)

var (
	errNumberExceedsDishes = errors.New("number exceeds available dishes")
	errInvalidType         = errors.New("invalid type")
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred", err)

			return
		}
	}()

	var amount int
	if _, err := fmt.Scan(&amount); err != nil {
		fmt.Println("cannot scan variable: ", err)

		return
	}

	dishes := &intheap.IntHeap{}

	for range amount {
		var pref int
		if _, err := fmt.Scan(&pref); err != nil {
			fmt.Println("cannot scan variable: ", err)

			return
		}

		heap.Push(dishes, pref)
	}

	var number int
	if _, err := fmt.Scan(&number); err != nil {
		fmt.Println("cannot scan variable: ", err)

		return
	}

	result, err := findDish(dishes, number)
	if err != nil {
		fmt.Println("cannot find preferred dish: ", err)
	}

	fmt.Println(result)
}

func findDish(dishes *intheap.IntHeap, number int) (int, error) {
	if number > dishes.Len() {
		return 0, fmt.Errorf("number %d more than amount of dishes %d: %w", number, dishes.Len(), errNumberExceedsDishes)
	}

	for range dishes.Len() - number {
		heap.Pop(dishes)
	}

	x := heap.Pop(dishes)

	value, ok := x.(int)
	if !ok {
		return 0, errInvalidType
	}

	return value, nil
}
