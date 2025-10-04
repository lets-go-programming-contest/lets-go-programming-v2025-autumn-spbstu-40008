package main

import (
	"fmt"
)

func main() {

	var (
		y, n, minT, maxT, temp int
		sign                   string
	)

	if _, err := fmt.Scan(&y); err != nil {
		return
	}

	for i := 0; i < y; i++ {
		if _, err := fmt.Scan(&n); err != nil {
			return
		}

		minT = 15
		maxT = 30
		for j := 0; j < n; j++ {
			if _, err := fmt.Scanf("%s %d", &sign, &temp); err != nil {
				return
			}
			switch sign {
			case "<=":
				maxT = min(maxT, temp)
			case ">=":
				minT = max(minT, temp)
			default:
				continue
			}

			if minT > maxT {
				fmt.Println(-1)
			} else {
				fmt.Println(minT)
			}
		}
	}
}
