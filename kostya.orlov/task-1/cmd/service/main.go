package main

import (
	"fmt"
	"os"
)

func main() {
	var (
		numb1    int
		numb2    int
		operator string
	)

	_, err := fmt.Scanln(&numb1)

	if err != nil {
		fmt.Println("Invalid first operand")
		os.Exit(0)
	}

	_, err = fmt.Scanln(&numb2)

	if err != nil {
		fmt.Println("Invalid second operand")
		os.Exit(0)
	}

	_, err = fmt.Scanln(&operator)

	if err != nil {
		fmt.Println("Invalid operation")
		os.Exit(0)
	}
	if operator != "+" &&
		operator != "-" &&
		operator != "/" &&
		operator != "*" {
		fmt.Println("Invalid operation")
		os.Exit(0)
	}

	if numb2 == 0 && operator == "/" {
		fmt.Println("Division by zero")
		os.Exit(0)
	}

	switch operator {
	case "+":
		fmt.Println(numb1 + numb2)
	case "-":
		fmt.Println(numb1 - numb2)
	case "*":
		fmt.Println(numb1 * numb2)
	case "/":
		fmt.Println(numb1 / numb2)
	}
}
