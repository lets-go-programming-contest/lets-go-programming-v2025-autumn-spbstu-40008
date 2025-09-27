package main

import "fmt"

func main() {
	var a, b int
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
	var op string
	_, err = fmt.Scan(&op)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}
	var result int
	if op == "+" {
		result = a + b
	} else if op == "-" {
		result = a - b
	} else if op == "*" {
		result = a * b
	} else if op == "/" {
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = a / b
	} else {
		fmt.Println("Invalid operation")
		return
	}
	fmt.Println(result)
}
