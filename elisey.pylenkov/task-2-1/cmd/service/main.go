package main

import "fmt"

func main() {
	var n, k int

	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	for i := 0; i < n; i++ {
		if _, err := fmt.Scan(&k); err != nil {
			return
		}

		minTemp, maxTemp := 15, 30

		for j := 0; j < k; j++ {
			var operator string
			var temp int

			if _, err := fmt.Scan(&operator, &temp); err != nil {
				return
			}

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
