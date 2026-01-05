package main

import (
	"fmt"
)

func main() {
	var first, second int
	var operator string
	_, err := fmt.Scan(&first)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&second)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&operator)
	if err != nil {
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
