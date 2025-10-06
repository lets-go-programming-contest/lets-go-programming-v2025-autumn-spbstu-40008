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
		return
	}

	for range countDepartments {
		if _, err := fmt.Scan(&countEmployees); err != nil {
			return
		}

		currentMin := MinTempDefault
		currentMax := MaxTempDefault

		for range countEmployees {
			if _, err := fmt.Scan(&setTempSign, &temp); err != nil {
				return
			}

			switch setTempSign {
			case ">=":
				currentMin = max(MinTempDefault, temp)
			case "<=":
				currentMax = min(MaxTempDefault, temp)
			default:
				continue
			}

			if currentMax < currentMin {
				fmt.Println(ErrVal)
			}
			fmt.Println(currentMin)
		}
	}
}
