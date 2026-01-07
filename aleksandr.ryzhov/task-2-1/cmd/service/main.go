package main

import "fmt"

func main() {
	var depCount, emplCount, minTemp, maxTemp, reqTemp int

	var reqTempInfo string

	_, err := fmt.Scan(&depCount)
	if err != nil || depCount < 1 {
		fmt.Println("Incorrect number of departments", err)

		return
	}

	for range depCount {
		_, err := fmt.Scan(&emplCount)
		if err != nil || emplCount < 1 {
			fmt.Println("Incorrect number of employees", err)

			return
		}

		maxTemp, minTemp = 30, 15

		for range emplCount {
			_, err := fmt.Scan(&reqTempInfo, &reqTemp)
			if err != nil {
				fmt.Println("Incorrect temperature information", err)

				return
			}

			switch reqTempInfo {
			case "<=":
				maxTemp = min(maxTemp, reqTemp)
			case ">=":
				minTemp = max(minTemp, reqTemp)
			default:
				fmt.Println("Incorrect temperature information")

				return
			}

			if minTemp > maxTemp {
				fmt.Println(-1)

				continue
			}

			fmt.Println(minTemp)
		}
	}
}
