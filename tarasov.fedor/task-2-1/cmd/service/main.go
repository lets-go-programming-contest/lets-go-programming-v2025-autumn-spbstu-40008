package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func processDepartment(scanner *bufio.Scanner) {
	if !scanner.Scan() {
		fmt.Println("Invalid number of employees")
		return
	}
	employeesCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number of employees")
		return
	}

	minTemp := 15
	maxTemp := 30
	fail := false

	for range employeesCount {
		if !scanner.Scan() {
			fmt.Println("Invalid temperature")
			os.Exit(0)
		}
		input := scanner.Text()

		var newTemp int
		var parseErr error

		switch {
		case strings.HasPrefix(input, ">="):
			numberStr := strings.TrimSpace(strings.TrimPrefix(input, ">="))
			newTemp, parseErr = strconv.Atoi(numberStr)
			if parseErr == nil {
				if newTemp > minTemp {
					minTemp = newTemp
				}
			} else {
				fmt.Println("Invalid number")
				os.Exit(0)
			}

		case strings.HasPrefix(input, "<="):
			numberStr := strings.TrimSpace(strings.TrimPrefix(input, "<="))
			newTemp, parseErr = strconv.Atoi(numberStr)
			if parseErr == nil {
				if newTemp < maxTemp {
					maxTemp = newTemp
				}
			} else {
				fmt.Println("Invalid number")
				os.Exit(0)
			}

		default:
			fmt.Println("Invalid operation")
			os.Exit(0)
		}

		if newTemp < 15 || newTemp > 30 || minTemp > maxTemp {
			fail = true
		}

		if !fail {
			fmt.Println(minTemp)
		}
	}

	if fail {
		fmt.Println(-1)
	}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		fmt.Println("Invalid number of departments")
		return
	}
	departmentsCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number of departments")
		return
	}

	for range departmentsCount {
		processDepartment(scanner)
	}
}
