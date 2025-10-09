package main

import "fmt"

const (
	MinTemp  = 15
	MaxTemp  = 30
	ErrorVal = -1
)

func main() {
	var (
		coutDepart, coutWorkers, Temp int
		Symbols                       string
	)

	_, err := fmt.Scan(&coutDepart)
	if err != nil {
		fmt.Println("Incorrect quantity departaments")

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
			_, err = fmt.Scan(&Symbols)
			if err != nil {
				fmt.Println("Incorrect symbol")

				return
			}

			_, err = fmt.Scan(&Temp)
			if err != nil {
				fmt.Println("Incorrect temperature")

				return
			}

			switch Symbols {
			case "<=":
				currentMax = min(currentMax, Temp)
			case ">=":
				currentMin = max(currentMin, Temp)
			default:
				continue
			}

			if currentMax < currentMin {
				fmt.Println(ErrorVal)
			} else {
				fmt.Println(currentMin)
			}
		}
	}
}
