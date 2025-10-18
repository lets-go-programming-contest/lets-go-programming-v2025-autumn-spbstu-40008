package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInputLine(input string) (string, int, error) {
	input = strings.TrimSpace(input)
	var op string
	var numStr string

	if strings.HasPrefix(input, ">=") {
		op = ">="
		numStr = strings.TrimSpace(strings.TrimPrefix(input, ">="))
	} else if strings.HasPrefix(input, "<=") {
		op = "<="
		numStr = strings.TrimSpace(strings.TrimPrefix(input, "<="))
	} else {
		return "", 0, fmt.Errorf("invalid operation")
	}

	val, err := strconv.Atoi(numStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid number")
	}

	return op, val, nil
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

	for i := 0; i < employeesCount; i++ {
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
		op, val, perr := parseInputLine(line)
		if perr != nil || val < 15 || val > 30 {
			fmt.Println(-1)
			fail = true
			continue
		}

		if op == ">=" {
			if val > minTemp {
				minTemp = val
			}
		} else { 
			if val < maxTemp {
				maxTemp = val
			}
		}

		if minTemp > maxTemp {
			fmt.Println(-1)
			fail = true
		} else {
			fmt.Println(minTemp)
		}
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

	for d := 0; d < departments; d++ {
		processDepartment(scanner)
	}
}
