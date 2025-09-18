package main

import (
	"fmt"
	"os"
)

func main() {
	var a, b int
	var operation string

	_, err := fmt.Scanln(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		os.Exit(0)
	}

	_, err = fmt.Scanln(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		os.Exit(0)
	}

	_, err = fmt.Scanln(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
		os.Exit(0)
	}

	var result int

	switch operation {
	case "+":
		result = a + b
		fmt.Println(result)

	case "-":
		result = a - b
		fmt.Println(result)

	case "*":
		result = a * b
		fmt.Println(result)

	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			os.Exit(0)
		}
		result = a / b
		fmt.Println(result)

	default:
		fmt.Println("Invalid operation")
		os.Exit(0)
	}

}
