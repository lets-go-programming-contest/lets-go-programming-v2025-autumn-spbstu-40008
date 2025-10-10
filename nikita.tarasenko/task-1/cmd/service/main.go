package main

import "fmt"

func main() {
	var operand1, operand2 int
	var operation string
	_, err := fmt.Scan(&operand1)

	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&operand2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
	}

	var result int
	switch operation {
	case "+":
		result = operand1 + operand2
	case "-":
		result = operand1 - operand2
	case "*":
		result = operand1 * operand2
	case "/":
		if operand2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = operand1 / operand2
	default:
		fmt.Println("Invalid operation")
	}

	fmt.Println(result)
}
