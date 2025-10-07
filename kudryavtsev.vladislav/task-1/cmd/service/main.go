package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var a int
	var b int
	var c string

	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		if _, err := bufio.NewReader(os.Stdin).ReadString('\n'); err != nil {
			fmt.Println("Error reading input:", err)
		}
		return
	}

	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		if _, err := bufio.NewReader(os.Stdin).ReadString('\n'); err != nil {
			fmt.Println("Error reading input:", err)
		}
		return
	}

	if _, err := fmt.Scan(&c); err != nil {
		fmt.Println("Invalid operation input:", err)
		return
	}

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
			return
		}
		fmt.Println(a / b)
	default:
		fmt.Println("Invalid operation")
	}
}
