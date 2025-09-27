package main

import "fmt"

func main() {
	var x, y int
	var operator string
	_, err := fmt.Scan(&x)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&y)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&operator)
	if err != nil {
		fmt.Println("Invalid input operation")
		return
	}

	switch operator {
	case "+":
		fmt.Println(x + y)
	case "-":
		fmt.Println(x - y)
	case "*":
		fmt.Println(x * y)
	case "/":
		if y == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(x / y)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
