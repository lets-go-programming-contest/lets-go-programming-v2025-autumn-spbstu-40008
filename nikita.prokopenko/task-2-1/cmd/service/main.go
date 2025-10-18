package main

import (
	"fmt"
)

func processDepartment(staffCount int) {
	mintemp := 15
	maxtemp := 30

	for i := 0; i < staffCount; i++ {
		var sign string
		var t int
		_, err := fmt.Scan(&sign, &t)
		if err != nil {
			return
		}

		if sign == "<=" {
			if t < maxtemp {
				maxtemp = t
			}
		} else if sign == ">=" {
			if t > mintemp {
				mintemp = t
			}
		}

		if mintemp > maxtemp {
			fmt.Println(-1)
		} else {
			fmt.Println(mintemp)
		}
	}
}

func main() {
	var n, k int
	_, err := fmt.Scan(&n)
	if err != nil {
		return
	}
	_, err = fmt.Scan(&k)
	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		processDepartment(k)
	}
}
