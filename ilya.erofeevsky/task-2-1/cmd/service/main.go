package main

import "fmt"

const (
	MinTemp  = 15
	MaxTemp  = 30
	ErrorVal = -1
)

func main() {
	var (
		coutDepart, coutWorkers, temp int
		symbols                       string
	)

	_, err := fmt.Scan(&coutDepart)
	if err != nil {
		fmt.Println("Error reading count of departaments")

		return
	}

	for range coutDepart {
		_, err = fmt.Scan(&coutWorkers)
		if err != nil {
			fmt.Println("Incorrect quantity workers")

			return
		}

		currentMax := MaxTemp
		currentMin := MinTemp

		for range coutWorkers {
			_, err = fmt.Scan(&symbols)
			if err != nil {
				fmt.Println("Incorrect symbol")

				return
			}

			_, err = fmt.Scan(&temp)
			if err != nil {
				fmt.Println("Incorrect temperature")

				return
			}

			switch symbols {
			case "<=":
				currentMax = min(currentMax, temp)
			case ">=":
				currentMin = max(currentMin, temp)
			default:
				fmt.Println("Warning: Unknow symbol")

				continue
			}

			if currentMax < currentMin {
				fmt.Println(ErrorVal)

				continue
			}
			
			fmt.Println(currentMin)
		}
	}
}
