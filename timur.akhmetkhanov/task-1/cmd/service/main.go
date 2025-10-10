package main

import (
	"fmt"
)

func main() {
	var (
		op1, op2  int
		operation string
	)

	_, err := fmt.Scan(&op1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&op2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
	}

	switch operation {
	case "+":
		fmt.Println(op1 + op2)
	case "-":
		fmt.Println(op1 - op2)
	case "*":
		fmt.Println(op1 * op2)
	case "/":
		if op2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(op1 / op2)
	default:
		fmt.Println("Invalid operation")
	}
}
