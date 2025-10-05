package main

import (
	"fmt"
)

func main() {

	var firstNum, secondNum int
	var operator string
	var result int

	_, err := fmt.Scanln(&firstNum)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scanln(&secondNum)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&operator)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {
	case "+":
		result = firstNum + secondNum
	case "-":
		result = firstNum - secondNum
	case "*":
		result = firstNum * secondNum
	case "/":
		if secondNum == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = firstNum / secondNum
	default:
		fmt.Println("Invalid operation")
		return
	}
	fmt.Println(result)
}
