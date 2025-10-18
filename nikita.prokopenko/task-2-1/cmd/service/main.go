go
package main

import (
	"fmt"
)

func main() {
	var n, k int

	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	if _, err := fmt.Scan(&k); err != nil {
		return
	}

	for range make([]struct{}, n) {
		min := 15
		max := 30

		for range make([]struct{}, k) {
			var sign string
			var t int

			if _, err := fmt.Scan(&sign, &t); err != nil {
				return
			}

			if sign == ">=" {
				if t > min {
					min = t
				}
			} else if sign == "<=" {
				if t < max {
					max = t
				}
			}

			if min > max {
				fmt.Println(-1)
			} else {
				fmt.Println(max)
			}
		}
	}
}
