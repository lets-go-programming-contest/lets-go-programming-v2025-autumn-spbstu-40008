package main

import "fmt"

func main() {
	var fstNum int
	var secNum int
	var op string

	_, err := fmt.Scan(&fstNum)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&secNum)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&op)
	if err != nil || !rightOp(op) {
		fmt.Println("Invalid operation")
		return
	}

	if op == "/" && secNum == 0 {
		fmt.Println("Division by zero")
		return
	}

	result := calculate(fstNum, secNum, op)
	fmt.Println(result)
}

func rightOp(op string) bool {
	switch op {
	case "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func calculate(a int, b int, op string) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	default:
		return 0
	}
}
