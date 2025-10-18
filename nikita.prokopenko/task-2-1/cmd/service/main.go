package main

import (
	"fmt"
)

func processDepartment(deptNum, staff int) {
	maxtemp := 30
	mintemp := 15
	for j := 1; j <= staff; j++ {
		fmt.Printf("Введите оператор и температуру (<= или >= число) сотрудник %d отдел %d:\n", j, deptNum)
		var temperatureData string
		var degrees int
		if _, err := fmt.Scan(&temperatureData, &degrees); err != nil {
			panic(err)
		}
		if degrees < 15 || degrees > 30 {
			panic("Температура вне допустимого диапазона")
		}
		if temperatureData != "<=" && temperatureData != ">=" {
			panic("Неверно введен оператор")
		}
		if temperatureData == "<=" && degrees < maxtemp {
			maxtemp = degrees
		} else if temperatureData == ">=" && degrees > mintemp {
			mintemp = degrees
		}
		if mintemp > maxtemp {
			fmt.Printf("Температура отдела %d после сотрудника %d: -1\n", deptNum, j)
		} else {
			fmt.Printf("Температура отдела %d после сотрудника %d: %d\n", deptNum, j, mintemp)
		}
	}
}

func main() {
	fmt.Println("Введите количество отделов:")
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		panic(err)
	}
	if departments < 1 || departments > 1000 {
		panic("Количество отделов вне диапазона")
	}
	fmt.Println("Введите количество сотрудников:")
	var staff int
	if _, err := fmt.Scan(&staff); err != nil {
		panic(err)
	}
	if staff < 1 || staff > 1000 {
		panic("Количество сотрудников вне диапазона")
	}
	for i := 1; i <= departments; i++ {
		processDepartment(i, staff)
	}
}
