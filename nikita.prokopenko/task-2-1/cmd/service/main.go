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
	ErrInvalidOperation      = errors.New("invalid operation")
	ErrInvalidNumber         = errors.New("invalid number")
	ErrReadEmployeesCount    = errors.New("could not read employees count")
	ErrParsingEmployeesCount = errors.New("failed to parse employees count")
	ErrProcessEmployeesData  = errors.New("minTemp exceeded maxTemp")
)

func parseInputLine(inputLine string) (string, int, error) {
	inputLine = strings.TrimSpace(inputLine)

	var operation string
	var numberStr string

	switch {
	case strings.HasPrefix(inputLine, ">="):
		operation = ">="
		numberStr = strings.TrimSpace(strings.TrimPrefix(inputLine, ">="))
	case strings.HasPrefix(inputLine, "<="):
		operation = "<="
		numberStr = strings.TrimSpace(strings.TrimPrefix(inputLine, "<="))
	default:
		return "", 0, ErrInvalidOperation
	}

	value, err := strconv.Atoi(numberStr)
	if err != nil {
		return "", 0, ErrInvalidNumber
	}

	return operation, value, nil
}

func parseEmployeesCount(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, ErrReadEmployeesCount
	}

	employeesCountStr := strings.TrimSpace(scanner.Text())

	employeesCount, err := strconv.Atoi(employeesCountStr)
	if err != nil {
		return 0, ErrParsingEmployeesCount
	}

	return employeesCount, nil
}

func processEmployeeData(line string, currentMin, currentMax int) (int, int, error) {
	operation, value, parseErr := parseInputLine(line)

	if parseErr != nil || value < 15 || value > 30 {
		return currentMin, currentMax, parseErr
	}

	switch operation {
	case ">=":
		if value > currentMin {
			currentMin = value
		}
	case "<=":
		if value < currentMax {
			currentMax = value
		}
	}

	if currentMin > currentMax {
		return currentMin, currentMax, ErrProcessEmployeesData
	}

	return currentMin, currentMax, nil
}

func processDepartment(scanner *bufio.Scanner) {
	employeesCount, err := parseEmployeesCount(scanner)
	if err != nil {
		fmt.Println(-1)

		return
	}

	minTemp := 15
	maxTemp := 30
	failed := false

	for i := 0; i < employeesCount; i++ {
		if failed || !scanner.Scan() {
			fmt.Println(-1)

			failed = true

			continue
		}

		line := scanner.Text()

		newMin, newMax, procErr := processEmployeeData(line, minTemp, maxTemp)

		if procErr != nil {
			fmt.Println(-1)

			failed = true

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

	for deptIndex := 0; deptIndex < departments; deptIndex++ {
		processDepartment(scanner)
	}
}
