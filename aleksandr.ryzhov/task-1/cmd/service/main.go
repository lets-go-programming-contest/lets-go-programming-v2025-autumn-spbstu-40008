package main

import "fmt"

func main() {
	var firstOperand, secondOperand int
	var operation string

	_, err := fmt.Scan(&firstOperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&secondOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(firstOperand + secondOperand)
	case "-":
		fmt.Println(firstOperand - secondOperand)
	case "*":
		fmt.Println(firstOperand * secondOperand)
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(firstOperand / secondOperand)
		}
	default:
		fmt.Println("Invalid operation")
	}
}
