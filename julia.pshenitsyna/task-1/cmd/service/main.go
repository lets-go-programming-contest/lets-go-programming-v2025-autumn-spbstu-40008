package main

import (
	"fmt"
)


func add(a float64, b float64) {
	fmt.Println(a + b)
}


func sub(a float64, b float64) {
	fmt.Println(a - b)
}


func mult(a float64, b float64) {
	fmt.Println(a * b)
}


func div(a float64, b float64) {
	if b == 0 {
		fmt.Println("Division by zero")
		return
	}
	fmt.Println(a / b)
}


func main() {

	var (
		a float64
		b float64
		op string
	)


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


	switch op {
	case "+":
		add(a, b)
	case "-":
		sub(a, b)
	case "*":
		mult(a, b)
	case "/":
		div(a, b)
	default:
		fmt.Println("Invalid operation")
	}
}
