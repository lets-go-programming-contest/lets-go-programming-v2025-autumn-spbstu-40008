package main

import "fmt"

func main() {
	var a, b int
	var symbol string
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Print("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Print("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&symbol)
	if err != nil {
		fmt.Print("Invalid operation")
		return
	}
	switch (symbol) {
		case "+":
			fmt.Print(a+b)
		case "-":
			fmt.Print(a-b)
		case "*":
			fmt.Print(a*b)
		case "/":
			if b == 0 {
				fmt.Print("Division by zero")
				return
			}
			fmt.Print(a/b)
		default:
			fmt.Print("Invalid operation")
	}
}