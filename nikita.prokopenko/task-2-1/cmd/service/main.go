package main

import (
	"fmt"
)

func processDepartment(staffCount int) {
	maxtemp := 30
	mintemp := 15

	for i := 0; i < staffCount; i++ {
		var sign string
		var temp int

		_, err := fmt.Scan(&sign, &temp)
		if err != nil {
			return
		}

		if sign == "<=" {
			if temp < maxtemp {
				maxtemp = temp
			}
		} else if sign == ">=" {
			if temp > mintemp {
				mintemp = temp
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
	fmt.Scan(&n)
	fmt.Scan(&k)

	for i := 0; i < n; i++ {
		processDepartment(k)
	}
}
