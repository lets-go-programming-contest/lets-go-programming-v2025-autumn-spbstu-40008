package main

import "fmt"

func conclusion(information any) {
	fmt.Println(information)
}

func mathematicalOperations(number1 int, number2 int, operation string) {
	switch operation {
	case "+":
		conclusion(number1 + number2)
	case "-":
		conclusion(number1 - number2)
	case "*":
		conclusion(number1 * number2)
	case "/":
		if number2 == 0 {
			conclusion("Division by zero")
		} else {
			conclusion(number1 / number2)
		}
	default:
		conclusion("Invalid operation")
	}
}

func main() {
	var number1, number2 int
	var operation string

	conclusion("Введите первое число:")
	if _, err := fmt.Scanln(&number1); err != nil {
		conclusion("Invalid first operand")
		return
	}

	conclusion("Введите второе число:")
	if _, err := fmt.Scanln(&number2); err != nil {
		conclusion("Invalid second operand")
		return
	}

	conclusion("Введите символ операции:")
	if _, err := fmt.Scanln(&operation); err != nil {
		conclusion("Invalid operation")
		return
	}

	mathematicalOperations(number1, number2, operation)
}
