package main

import (
	"fmt"
	"os"
)

func main() {

	var (
		n int
	)

	if _, err := fmt.Scan(&n); err != nil {
		fmt.Println("Error input N")
		os.Exit(0)
	}

	for i := 0; i < n; i++ {
		var (
			k, temp  int
			upper    = 30
			lower    = 15
			operator string
		)

		if _, err := fmt.Scan(&k); err != nil {
			fmt.Println("Error input K")
			os.Exit(0)
		}

		for j := 0; j < k; j++ {
			if _, err := fmt.Scan(&operator, &temp); err != nil {
				fmt.Println("Error input temperature")
				os.Exit(0)
			}

			switch operator {
			case "<=":
				upper = min(upper, temp)
			case ">=":
				lower = max(lower, temp)
			default:
				fmt.Println("Invalid operator")
				os.Exit(0)
			}

			if lower > upper {
				fmt.Println(-1)
			} else {
				fmt.Println(lower)
			}
		}
	}

}
