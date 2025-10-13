package main

import "fmt"

func main() {
	var (
		departments, employees, minT, maxT, temp int
		sign                                     string
	)

	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println("Invalid departments value")
		return
	}

	for range departments {
		if _, err := fmt.Scan(&employees); err != nil {
			fmt.Println("Invalid employees value")
			return
		}

		minT, maxT = 15, 30

		for range employees {
			if _, err := fmt.Scanf("%s %d", &sign, &temp); err != nil {
				fmt.Println("Invalid temp format")
				return
			}

			switch sign {
			case "<=":
				maxT = min(maxT, temp)
			case ">=":
				minT = max(minT, temp)
			default:
				fmt.Println("Invalid sign")
			}

			if minT > maxT {
				fmt.Println(-1)

				continue
			}

			fmt.Println(minT)
		}
	}
}
