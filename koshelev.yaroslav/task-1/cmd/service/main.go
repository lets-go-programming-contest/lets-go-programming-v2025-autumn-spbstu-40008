package main

import "fmt"

func main() {
	var fstNum int
	var secNum int
	var op string

	_, err1 := fmt.Scan(&fstNum)
	if err1 != nil {
		fmt.Println("Invalid first operand:", err1)
		return
	}
	
	_, err2 := fmt.Scan(&secNum)
	if err2 != nil {
		fmt.Println("Invalid second operand:", err2)
		return
	}

	_, _ = fmt.Scan(&op)

	if !rightOp(op) {
		fmt.Println("Invalid operation")
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
		if b == 0 {
			fmt.Println("Division by zero")
			return 0
		}
		return a / b
	default:
		return 0
	}
}
