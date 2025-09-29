package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
	ErrInvalidNumber    = errors.New("invalid number")
)

func parseInputLine(input string) (string, int, error) {
	var operation string

	var numStr string

	switch {
	case strings.HasPrefix(input, ">="):
		operation = ">="
		numStr = strings.TrimSpace(strings.TrimPrefix(input, ">="))
	case strings.HasPrefix(input, "<="):
		operation = "<="
		numStr = strings.TrimSpace(strings.TrimPrefix(input, "<="))
	default:

		return "", 0, ErrInvalidOperation
	}

	val, err := strconv.Atoi(numStr)
	if err != nil {

		return "", 0, ErrInvalidNumber
	}

	return operation, val, nil
}

func parseEmployeesCount(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, fmt.Errorf("could not read employees count")
	}

	return strconv.Atoi(scanner.Text())
}

func processEmployeeData(input string, minTemp, maxTemp int) (int, int, error) {
	operation, val, err := parseInputLine(input)
	if err != nil || val < 15 || val > 30 {
		return minTemp, maxTemp, fmt.Errorf("invalid input line")
	}

	switch operation {
	case ">=":
		if val > minTemp {
			minTemp = val
		}
	case "<=":
		if val < maxTemp {
			maxTemp = val
		}
	}

	if minTemp > maxTemp {
		return minTemp, maxTemp, fmt.Errorf("minTemp exceeded maxTemp")
	}

	return minTemp, maxTemp, nil
}

func processDepartment(scanner *bufio.Scanner) {
	employeesCount, err := parseEmployeesCount(scanner)
	if err != nil {
		fmt.Println(-1)
		return
	}

	minTemp := 15
	maxTemp := 30
	fail := false

	for range employeesCount {
		if !scanner.Scan() {
			fail = true
			break
		}

		input := scanner.Text()
		newMinTemp, newMaxTemp, err := processEmployeeData(input, minTemp, maxTemp)
		if err != nil {
			fail = true
			continue
		}

		minTemp = newMinTemp
		maxTemp = newMaxTemp
		fmt.Println(minTemp)
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
