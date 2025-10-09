package main

import (
	"fmt"
)

func main() {
	var (
		N, K, curTemp, minTemp, maxTemp int
		sign                            string
	)

	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Println("Invalid value")
	}

	for range N {
		_, err = fmt.Scan(&K)
		if err != nil {
			fmt.Println("Invalid value")
		}
		maxTemp = 30
		minTemp = 15

		for range K {
			fmt.Scan(&sign, &curTemp)
			if curTemp <= 30 && curTemp >= 15 {
				switch sign {
				case "<=":
					maxTemp = min(maxTemp, curTemp)
				case ">=":
					minTemp = max(minTemp, curTemp)
				default:
					fmt.Println("Invalid operation")
				}
			} else {
				fmt.Println(-1)
			}
			if maxTemp < minTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(min(minTemp, maxTemp))
			}
		}
	}
}
