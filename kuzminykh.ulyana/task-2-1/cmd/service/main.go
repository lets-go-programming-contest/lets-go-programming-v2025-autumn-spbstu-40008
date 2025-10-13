package main

import (
	"fmt"
)

func main() {
	var (
		dep, emp, curTemp, minTemp, maxTemp int
		sign                                string
	)

	_, err := fmt.Scan(&dep)
	if err != nil {
		fmt.Println(err)

		return
	}

	for range dep {
		_, err = fmt.Scan(&emp)
		if err != nil {
			fmt.Println(err)

			return
		}

		maxTemp = 30
		minTemp = 15

		for range emp {
			_, err = fmt.Scan(&sign, &curTemp)
			if err != nil {
				fmt.Println(err)

				return
			}

			switch sign {
			case "<=":
				maxTemp = min(maxTemp, curTemp)
			case ">=":
				minTemp = max(minTemp, curTemp)
			default:
				fmt.Println("Invalid operation")

				return
			}

			if maxTemp < minTemp {
				fmt.Println(-1)

				continue
			}

			fmt.Println(min(minTemp, maxTemp))
		}
	}
}
