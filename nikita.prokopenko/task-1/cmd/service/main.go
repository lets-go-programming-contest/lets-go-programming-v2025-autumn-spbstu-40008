package main

import "fmt"

func conclusion(information any){
	fmt.Println(information)
}


func mathematical_operations(number1 int, number2 int, the_addition string){
	switch the_addition{
	case "+":
		conclusion(number1 + number2)
	case "-":	
		conclusion(number1 - number2)
	case "*":
		conclusion(number1 * number2)
	case "/":
		if number2 == 0{
			conclusion("Division by zero")
		}else{
			conclusion(number1 / number2)
		}
	default:
		fmt.Println("Invalid operation")
	}
}

func main(){
	var number1, number2 int
	var the_addition string
	conclusion("Введите первое число:")
	if _,err := fmt.Scanln(&number1); err != nil{
		conclusion("Invalid first operand")
		return
	}
	conclusion("Введите второе число:")
	if _,err := fmt.Scanln(&number2); err != nil{
		conclusion("Invalid second operand")
		return
	}
	conclusion("Введите символ операции:")
	if _,err := fmt.Scanln(&the_addition); err != nil{
		conclusion("Invalid operation")
		return
	}
	mathematical_operations(number1, number2,the_addition)

} 
