package main

import "fmt"

const (
	MinTempDefault = 15
	MaxTempDefault = 30
	ErrVal         = -1
)

func main() {
	var (
		countDepartments, countEmployees, temp int
		setTempSign                            string
	)

	if _, err := fmt.Scan(&countDepartments); err != nil {
		fmt.Printf("Bad input for departments: %v\n", err)

		return
	}

	for range countDepartments {
		if _, err := fmt.Scan(&countEmployees); err != nil {
			fmt.Printf("Bad input for employees: %v\n", err)

			return
		}

		currentMin := MinTempDefault
		currentMax := MaxTempDefault

		for range countEmployees {
			if _, err := fmt.Scan(&setTempSign, &temp); err != nil {
				fmt.Printf("Bad input: %v\n", err)

				return
			}

			if setTempSign != "<=" && setTempSign != ">=" {
				fmt.Printf("Invalid sign: %s\n", setTempSign)

				return
			}

			switch setTempSign {
			case "<=":
				currentMax = min(currentMax, temp)
			case ">=":
				currentMin = max(currentMin, temp)
			default:
				continue
			}

			if currentMax < currentMin {
				fmt.Println(ErrVal)
			} else {
				fmt.Println(currentMin)
			}
		}
	}
}
