package main

import (
	"fmt"
	"time"
)

func add(a float64, b float64) {
	time.Sleep(1 * time.Second)
	fmt.Print(a + b)
}

func sub(a float64, b float64) {
	time.Sleep(1 * time.Second)
	fmt.Print(a - b)
}

func mult(a float64, b float64) {
	time.Sleep(1 * time.Second)
	fmt.Print(a * b)
}

func div(a float64, b float64) {
	if b == 0 {
		fmt.Println("Division by zero")
		return
	}
	time.Sleep(1 * time.Second)
	fmt.Print(a / b)
}

func main() {
	var a float64
	var b float64
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
	if err != nil{
		fmt.Println("Invalid operation")
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
