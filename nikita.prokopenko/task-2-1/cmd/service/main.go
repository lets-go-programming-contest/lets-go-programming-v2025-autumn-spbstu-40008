package main

import (
	"fmt"
)

func processDepartment(deptNum, staffCount int) {
	maxtemp := 30
	mintemp := 15
	fail := false

	for employeeIndex := 1; employeeIndex <= staffCount; employeeIndex++ {
		if fail {
			fmt.Println(-1)
			continue
		}

		var temperatureData string
		var degrees int

		if _, err := fmt.Scan(&temperatureData, &degrees); err != nil {
			fmt.Println(-1)
			fail = true
			continue
		}

		if degrees < 15 || degrees > 30 {
			fmt.Println(-1)
			fail = true
			continue
		}

		if temperatureData != "<=" && temperatureData != ">=" {
			fmt.Println(-1)
			fail = true
			continue
		}

		if temperatureData == "<=" && degrees < maxtemp {
			maxtemp = degrees
		} else if temperatureData == ">=" && degrees > mintemp {
			mintemp = degrees
		}

		if mintemp > maxtemp {
			fmt.Println(-1)
			fail = true
		} else {
			fmt.Println(mintemp)
		}
	}
}

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		return
	}
	if departments < 1 || departments > 1000 {
		return
	}

	for deptNum := 1; deptNum <= departments; deptNum++ {
		var staffCount int
		if _, err := fmt.Scan(&staffCount); err != nil {
			return
		}
		if staffCount < 1 || staffCount > 1000 {
			return
		}

		processDepartment(deptNum, staffCount)
	}
}
