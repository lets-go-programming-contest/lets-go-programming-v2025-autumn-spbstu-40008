package main

import (
	"fmt"
)

func main() {
	var (
		departments, employees, maxT, minT, temperature int
		sign                                            string
	)
	fmt.Println("Введите отделы:")
	_, err := fmt.Scan(&departments)
	if err != nil {
		return
	}
	for range departments { // считаем по каждому отделу

		fmt.Println("Введите сотрудников:")
		_, err := fmt.Scan(&employees)
		if err != nil {
			return
		}
		maxT = 10000000000000
		minT = 0
		for range employees { //считаем по каждому сторуднику

			fmt.Println("Введите температуру:")
			_, err := fmt.Scanf("\n%s %d", &sign, &temperature)
			if err != nil {
				return
			}
			switch sign {
			case ">=":
				minT = max(minT, temperature)
			case "<=":
				maxT = min(maxT, temperature)
			default:
				continue
			}
			if minT > maxT{
				fmt.Println(-1)
				continue
			}
			fmt.Println(min(minT, maxT))
		}
	}
}
