package main

import (
	"fmt"
)

func processDepartment(staffCount int) {
	mintemp := 15
	maxtemp := 30

	for i := 0; i < staffCount; i++ {
		var op string
		var val int

		if _, err := fmt.Scan(&op, &val); err != nil {
			return
		}

		if op == "<=" {
			if val < maxtemp {
				maxtemp = val
			}
		} else if op == ">=" {
			if val > mintemp {
				mintemp = val
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
	var departments, staffCount int

	if _, err := fmt.Scan(&departments); err != nil {
		return
	}
	if _, err := fmt.Scan(&staffCount); err != nil {
		return
	}

	for i := 0; i < departments; i++ {
		processDepartment(staffCount)
	}
}
