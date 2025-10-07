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
		fmt.Print("Invalid first operand\n")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}
	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Print("Invalid second operand\n")
		bufio.NewReader(os.Stdin).ReadString('\n')
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
			return
		}
		fmt.Println(a / b)
	default:
		fmt.Println("Invalid operation")
	}
}
