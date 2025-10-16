package main

import (
	"fmt"
)

const (
	errorValue          = -1
	defaultMinTempValue = 15
	defaultMaxTempValue = 30
)

func main() {
	var (
		departmentCount, countEmployess, temp int
		wishTempSign                          string
	)

	if _, error := fmt.Scan(&departmentCount); error != nil {
		fmt.Printf("Bad input for departments: %v\n", error)
		return
	}

	for range departmentCount {
		if _, error := fmt.Scan(&countEmployess); error != nil {
			fmt.Printf("Bad input for employess: %v\n", error)
			return
		}
		var (
			departmentMinTempValue = defaultMinTempValue
			departmentMaxTempValue = defaultMaxTempValue
		)
		for range countEmployess {
			if _, error := fmt.Scan(&wishTempSign, &temp); error != nil {
				fmt.Printf("Bad input: %v\n", error)
				return
			}

			switch wishTempSign {
			case "<=":
				departmentMaxTempValue = min(departmentMaxTempValue, temp)
			case ">=":
				departmentMinTempValue = max(departmentMinTempValue, temp)
			default:
				fmt.Printf("Invalid sign: %s\n", wishTempSign)
				return
			}

			if departmentMaxTempValue < departmentMinTempValue {
				fmt.Println(errorValue)
			} else {
				fmt.Println(departmentMinTempValue)
			}

		}
	}

}
