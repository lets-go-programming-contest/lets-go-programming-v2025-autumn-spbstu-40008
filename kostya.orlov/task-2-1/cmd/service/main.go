package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readLine() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		return scanner.Text(), nil
	}

	err := scanner.Err()

	if err != nil {
		return "", err
	}

	return "", nil
}

func main() {
	var (
		n, k, upper, lower int
		comfortable        []int
		temp               string
	)

	_, err := fmt.Scan(&n)

	if err != nil {
		fmt.Println("Error input")
		os.Exit(0)
	}

	for i := 0; i < n; i++ {
		_, err = fmt.Scan(&k)

		if err != nil {
			fmt.Println("Error input")
			os.Exit(0)
		}

		upper = 30
		lower = 15

		for j := 0; j < k; j++ {
			temp, err = readLine()

			if err != nil {
				fmt.Println("Error input")
				os.Exit(0)
			}

			degree, err := strconv.Atoi(temp[3:])

			if err != nil {
				fmt.Println("Error converting")
				os.Exit(0)
			}

			temp = temp[:2]

			switch temp {
			case "<=":
				if degree <= upper {
					upper = degree
				}
			case ">=":
				if degree >= lower {
					lower = degree
				}
			}

			if lower > upper {
				comfortable = append(comfortable, -1)
			} else {
				comfortable = append(comfortable, lower)
			}

		}
	}

	for _, val := range comfortable {
		fmt.Println(val)
	}
}
