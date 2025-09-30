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

	_, err1 := fmt.Scanln(&numb1)

	if err1 != nil {
		fmt.Print("Invalid first operand\n")
		os.Exit(1)
	}

	_, err2 := fmt.Scanln(&numb2)

	if err2 != nil {
		fmt.Print("Invalid second operand\n")
		os.Exit(1)
	}

	_, err3 := fmt.Scanln(&operator)

	if err3 != nil {
		fmt.Print("Invalid operation\n")
		os.Exit(1)
	} else if operator != "+" &&
		operator != "-" &&
		operator != "/" &&
		operator != "*" {
		fmt.Print("Invalid operation\n")
		os.Exit(1)
	}

	if numb2 == 0 && operator == "/" {
		fmt.Print("Division by zero\n")
		os.Exit(1)
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
