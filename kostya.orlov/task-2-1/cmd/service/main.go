package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func readLine() (string, error) {
	line, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func main() {
	var (
		n, k, upper, lower int
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

			if err != nil || temp == "" {
				fmt.Println("Error input")
				os.Exit(0)
			}

			parts := strings.Fields(temp)

			if len(parts) != 2 {
				fmt.Println("Error")
				os.Exit(0)
			}

			operator := parts[0]
			degreeStr := parts[1]

			degree, err := strconv.Atoi(degreeStr)

			if err != nil {
				fmt.Println("Error converting")
				os.Exit(0)
			}

			switch operator {
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
				fmt.Println(-1)
			} else {
				fmt.Println(lower)
			}

		}
	}
}
