package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	initialMinTemp = 15
	initialMaxTemp = 30
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	defer func() {
		_ = writer.Flush()
	}()

	var departmentCount int
	if _, err := fmt.Fscan(reader, &departmentCount); err != nil {
		return
	}

	for range departmentCount {
		var employeeCount int
		if _, err := fmt.Fscan(reader, &employeeCount); err != nil {
			break
		}

		minTemp := initialMinTemp
		maxTemp := initialMaxTemp

		for range employeeCount {
			var (
				operator    string
				temperature int
			)

			if _, err := fmt.Fscan(reader, &operator, &temperature); err != nil {
				break
			}

			if operator == ">=" {
				minTemp = max(minTemp, temperature)
			} else {
				maxTemp = min(maxTemp, temperature)
			}

			if minTemp <= maxTemp {
				_, _ = fmt.Fprintln(writer, minTemp)
			} else {
				_, _ = fmt.Fprintln(writer, -1)
			}
		}
	}
}
