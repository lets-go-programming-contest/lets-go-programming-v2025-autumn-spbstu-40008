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

	if _, err := fmt.Scan(&departmentCount); err != nil {
		fmt.Printf("Bad input for departments: %v\n", err)

		return
	}

	for range departmentCount {
		if _, err := fmt.Scan(&countEmployess); err != nil {
			fmt.Printf("Bad input for employess: %v\n", err)

			return
		}
		var (
			departmentMinTempValue = defaultMinTempValue
			departmentMaxTempValue = defaultMaxTempValue
		)
		for range countEmployess {
			if _, err := fmt.Scan(&wishTempSign, &temp); err != nil {
				fmt.Printf("Bad input: %v\n", err)

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
