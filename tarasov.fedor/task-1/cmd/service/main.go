package main

import (
	"fmt"
	"os"
)

func main() {
	var (
		a, b      int
		operation string
	)

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

	switch operation {
	case "+":
		fmt.Println(a + b)

	case "-":
		fmt.Println(a - b)

	case "*":
		fmt.Println(a * b)

	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			os.Exit(0)
		}
		fmt.Println(a / b)

	default:
		fmt.Println("Invalid operation")
		os.Exit(0)
	}
}
