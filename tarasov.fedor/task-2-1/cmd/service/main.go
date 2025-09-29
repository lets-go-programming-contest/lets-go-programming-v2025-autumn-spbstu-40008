package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInputLine(input string) (string, int, error) {
	var op, numStr string

	if strings.HasPrefix(input, ">=") {
		op = ">="
		numStr = strings.TrimSpace(strings.TrimPrefix(input, ">="))
	} else if strings.HasPrefix(input, "<=") {
		op = "<="
		numStr = strings.TrimSpace(strings.TrimPrefix(input, "<="))
	} else {
		return "", 0, errors.New("invalid operation")
	}

	val, err := strconv.Atoi(numStr)
	if err != nil {
		return "", 0, errors.New("invalid number")
	}

	return op, val, nil
}

func processDepartment(scanner *bufio.Scanner) {
	if !scanner.Scan() {
		fmt.Println(-1)

		return
	}
	employeesCountStr := scanner.Text()
	employeesCount, err := strconv.Atoi(employeesCountStr)
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
		op, val, err := parseInputLine(input)
		if err != nil || val < 15 || val > 30 {
			fail = true
			continue
		}

		switch op {
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
			fail = true
			continue
		}

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
