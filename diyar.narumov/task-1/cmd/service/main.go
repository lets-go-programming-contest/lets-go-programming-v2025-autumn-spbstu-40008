package main

import "fmt"

func main() {
	var a, b int
	var op string

	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&op)
	if err != nil {
		fmt.Println("Invalid operation")
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
