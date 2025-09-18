package main

import (
	"fmt"
)

func main() {
	var a, b int
	var c string

	if _, err := fmt.Scan(&a); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scan(&b); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Scan(&c)
	switch c {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(a / b)
		}
	default:
		fmt.Println("Invalid operation")
	}

}
