package main

import "fmt"

func main() {
	var a, b int
	var op string

	if _, err := fmt.Scan(&a); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scan(&op); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	if _, err := fmt.Scan(&b); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var result int
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = a / b
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(result)
}
