package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		n int
		k int
	)

	_, err := fmt.Scanln(&n)
	if err != nil {
		fmt.Println("Invalid number of departments")
		os.Exit(0)
	}

	for i := 0; i < n; i++ {
		_, err := fmt.Scanln(&k)
		if err != nil {
			fmt.Println("Invalid number of employees")
			os.Exit(0)
		}

		var (
			optTemp  = 0
			highTemp = 30
			fail     = false
		)

		for j := 0; j < k; j++ {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Invalid temperature")
				os.Exit(0)
			}

			switch {
			case strings.HasPrefix(input, ">="):
				var numberStr = strings.TrimSpace(strings.TrimPrefix(input, ">="))
				number, err := strconv.Atoi(numberStr)
				if err != nil {
					fmt.Println("Invalid number")
					os.Exit(0)
				}
				if number < 15 || number > 30 {
					fail = true
					break
				}

				if optTemp == 0 || number > optTemp && number <= highTemp {
					optTemp = number
					fmt.Println(optTemp)
				} else {
					fail = true
					break
				}

			case strings.HasPrefix(input, "<="):
				var numberStr = strings.TrimSpace(strings.TrimPrefix(input, "<="))
				number, err := strconv.Atoi(numberStr)
				if err != nil {
					fmt.Println("Invalid number")
					os.Exit(0)
				}
				if number < 15 || number > 30 {
					fail = true
					break
				}

				if highTemp == 30 || number < highTemp && number >= optTemp {
					highTemp = number
					fmt.Println(optTemp)
				} else {
					fmt.Println(optTemp)
				}
			default:
				fmt.Println("Invalid operation")
				os.Exit(0)
			}
			if fail {
				fmt.Println(-1)
				break
			}
		}
	}
}
