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
	errInvalidOperation = errors.New("invalid operation")
	errInvalidNumber    = errors.New("invalid number")
)

func parseInputLine(input string) (string, int, error) {
	input = strings.TrimSpace(input)

	var operation, numberStr string

	switch {
	case strings.HasPrefix(input, ">="):
		operation = ">="
		numberStr = strings.TrimSpace(strings.TrimPrefix(input, ">="))
	case strings.HasPrefix(input, "<="):
		operation = "<="
		numberStr = strings.TrimSpace(strings.TrimPrefix(input, "<="))
	default:
		return "", 0, errInvalidOperation
	}

	value, err := strconv.Atoi(numberStr)
	if err != nil {
		return "", 0, errInvalidNumber
	}

	return operation, value, nil
}

func processSingleEmployeeLine(line string, currentMin, currentMax int) (int, int, bool) {
	operation, value, parseErr := parseInputLine(line)
	if parseErr != nil || value < 15 || value > 30 {
		return currentMin, currentMax, true
	}

	if operation == ">=" {
		if value > currentMin {
			currentMin = value
		}
	} else {
		if value < currentMax {
			currentMax = value
		}
	}

	if currentMin > currentMax {
		return currentMin, currentMax, true
	}

	return currentMin, currentMax, false
}

func processDepartment(scanner *bufio.Scanner) {
	if !scanner.Scan() {
		return
	}

	employeesCountStr := strings.TrimSpace(scanner.Text())

	employeesCount, err := strconv.Atoi(employeesCountStr)
	if err != nil {
		fmt.Println(-1)

		return
	}

	minTemp := 15
	maxTemp := 30
	fail := false

	for range make([]struct{}, employeesCount) {
		if !scanner.Scan() {
			fmt.Println(-1)

			fail = true

			continue
		}

		if fail {
			fmt.Println(-1)

			continue
		}

		line := scanner.Text()

		newMin, newMax, invalid := processSingleEmployeeLine(line, minTemp, maxTemp)
		if invalid {
			fmt.Println(-1)

			fail = true

			continue
		}

		minTemp = newMin
		maxTemp = newMax

		fmt.Println(minTemp)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		return
	}

	departmentsStr := strings.TrimSpace(scanner.Text())

	departments, err := strconv.Atoi(departmentsStr)
	if err != nil {
		return
	}

	for range make([]struct{}, departments) {
		processDepartment(scanner)
	}
}
