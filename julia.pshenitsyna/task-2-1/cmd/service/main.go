package main

import (
	"fmt"
)

func main() {
	var (
		departments, employees, maxT, minT, temperature int
		sign                                            string
	)

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("Input error")

		return
	}

	for range departments {
		_, err := fmt.Scan(&employees)
		if err != nil {
			fmt.Println("Input error")

			return
		}

		maxT = 30
		minT = 15

		for range employees {
			_, err := fmt.Scanf("%s %d", &sign, &temperature)
			if err != nil {
				fmt.Println("Input error")

				return
			}

			switch sign {
			case ">=":
				minT = max(minT, temperature)
			case "<=":
				maxT = min(maxT, temperature)
			default:
				fmt.Println("Wrong operator")
			}

			if minT > maxT {
				fmt.Println(-1)

				continue
			}

			fmt.Println(minT)
		}
	}
}
