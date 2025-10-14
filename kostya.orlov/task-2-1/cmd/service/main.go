package main

import (
	"fmt"
)

func main() {
	var (
		n, k, wish  int
		lower       = 15
		upper       = 30
		comfortable []int
		temp        string
	)

	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		fmt.Scan(&k)
		for j := 0; j < k; j++ {
			fmt.Scan(&temp)
			fmt.Scan(&wish)

			if temp == "<=" &&
				wish >= lower {
				upper = wish
				comfortable = append(comfortable, lower)
			} else if temp == ">=" &&
				wish <= upper {
				lower = wish
				comfortable = append(comfortable, lower)
			} else {
				comfortable = append(comfortable, -1)
			}
		}
		lower = 15
		upper = 30
	}

	for i := 0; i < len(comfortable); i++ {
		fmt.Println(comfortable[i])
	}
}
