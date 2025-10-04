package main

import "fmt"

func main() {
	var (
		departments, employees, minT, maxT, temp int
		sign                                     string
	)

	if _, err := fmt.Scan(&departments); err != nil {
		return
	}

	for range departments {
		if _, err := fmt.Scan(&employees); err != nil {
			return
		}

		minT, maxT = 15, 30

		for range employees {
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
