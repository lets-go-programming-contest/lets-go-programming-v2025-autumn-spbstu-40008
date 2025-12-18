package main

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	MinTemperature = 15
	MaxTemperature = 30
)

var (
	ErrInvalidOperator    = errors.New("invalid operator")
	ErrInvalidTemperature = errors.New("invalid temperature")
)

type TemperatureRange struct {
	min int
	max int
}

func NewTemperatureRange(minTemp, maxTemp int) *TemperatureRange {
	return &TemperatureRange{
		min: minTemp,
		max: maxTemp,
	}
}

func (tr *TemperatureRange) Update(operation string, temperature int) error {
	switch operation {
	case ">=":
		tr.min = max(tr.min, temperature)
	case "<=":
		tr.max = min(tr.max, temperature)
	default:
		return ErrInvalidOperator
	}

	return nil
}

func (tr *TemperatureRange) GetOptimalTemperature() int {
	if tr.min > tr.max {
		return -1
	}

	return tr.min
}

func ParseConstraint(oper string, tempStr string) (string, int, error) {
	if oper != "<=" && oper != ">=" {
		return "", 0, ErrInvalidOperator
	}

	tempInt, err := strconv.Atoi(tempStr)
	if err != nil {
		return "", 0, ErrInvalidTemperature
	}

	if tempInt > MaxTemperature || tempInt < MinTemperature {
		return "", 0, ErrInvalidTemperature
	}

	return oper, tempInt, nil
}

func processDepartment(employeeCount int) {
	temperatureRange := NewTemperatureRange(MinTemperature, MaxTemperature)

	for range employeeCount {
		var oper, tempStr string
		if _, err := fmt.Scan(&oper, &tempStr); err != nil {
			fmt.Printf("Error reading input: %v\n", err)

			continue
		}

		operation, temperature, err := ParseConstraint(oper, tempStr)
		if err != nil {
			fmt.Printf("Error parsing constraint: %v\n", err)

			continue
		}

		if err := temperatureRange.Update(operation, temperature); err != nil {
			fmt.Printf("Error updating temperature range: %v\n", err)

			continue
		}

		fmt.Println(temperatureRange.GetOptimalTemperature())
	}
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Printf("Error reading departments: %v\n", err)

		return
	}

	for range departments {
		var employees int

		_, err := fmt.Scan(&employees)
		if err != nil {
			fmt.Printf("Error reading employees: %v\n", err)

			return
		}

		processDepartment(employees)
	}
}
