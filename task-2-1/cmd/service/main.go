package main

import (
    "fmt"
)

func main() {
    fmt.Println("Введите количество отделов:")
    var departments int
    fmt.Scan(&departments)
    if departments < 1 || departments > 1000 {
        panic("Количество отделов вне диапазона")
    }

    fmt.Println("Введите количество сотрудников:")
    var staff int
    fmt.Scan(&staff)
    if staff < 1 || staff > 1000 {
        panic("Количество сотрудников вне диапазона")
    }

    for i := 1; i <= departments; i++ {
        maxtemp := 30
        mintemp := 15

        for j := 1; j <= staff; j++ {
            fmt.Println("Введите оператор и температуру (<= или >= число) сотрудник", j, "отдел", i, ":")
            var temperature_data string
            var degrees int
            fmt.Scan(&temperature_data, &degrees)

            if degrees < 15 || degrees > 30 {
                panic("Температура вне допустимого диапазона")
            }
            if temperature_data != "<=" && temperature_data != ">=" {
                panic("Неверно введен оператор")
            }

            if temperature_data == "<=" && degrees < maxtemp {
                maxtemp = degrees
            } else if temperature_data == ">=" && degrees > mintemp {
                mintemp = degrees
            }

            if mintemp > maxtemp {
                fmt.Println("Температура отдела", i, "после сотрудника", j, ":", -1)
            } else {
                fmt.Println("Температура отдела", i, "после сотрудника", j, ":", mintemp)
            }
        }
    }
}
