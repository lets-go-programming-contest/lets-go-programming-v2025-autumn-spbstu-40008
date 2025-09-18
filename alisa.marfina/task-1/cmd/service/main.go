package main

import (
	"fmt"
)

func main() {
	var first int
	var second int
	var operator string
	_, err1 := fmt.Scan(&first)
	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err2 := fmt.Scan(&second)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err3 := fmt.Scan(&operator)
	if err3 != nil {
		fmt.Println("Invalid input operation")
		return
	}

	switch operator {
	case "+":
		fmt.Println(first + second)
	case "-":
		fmt.Println(first - second)
	case "*":
		fmt.Println(first * second)
	case "/":
		if second == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(first / second)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
