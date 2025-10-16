package main

import "fmt"

func main() {
	var n, k int

	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		fmt.Scan(&k)

		minTemp, maxTemp := 15, 30

		for j := 0; j < k; j++ {
			var operator string
			var temp int

			fmt.Scan(&operator, &temp)

			switch operator {
			case ">=":
				if temp > minTemp {
					minTemp = temp
				}
			case "<=":
				if temp < maxTemp {
					maxTemp = temp
				}
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
