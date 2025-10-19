package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidComparisonOperator = errors.New("invalid comparison operator; use '<=' or '>='")

func processEmployee(low, high *int) error {
	var cmp string

	if _, err := fmt.Scan(&cmp); err != nil {
		return fmt.Errorf("failed to read comparison operator: %w", err)
	}

	var temperature int

	if _, err := fmt.Scan(&temperature); err != nil {
		return fmt.Errorf("failed to read temperature value: %w", err)
	}

	switch cmp {
	case "<=":
		if *high > temperature {
			*high = temperature
		}
	case ">=":
		if *low < temperature {
			*low = temperature
		}
	default:
		return ErrInvalidComparisonOperator
	}

	if *low > *high {
		fmt.Println(-1)
	} else {
		fmt.Println(*low)
	}

	return nil
}

func processDepartment() error {
	var employeesNumber int

	if _, err := fmt.Scan(&employeesNumber); err != nil {
		return fmt.Errorf("failed to read number of employees: %w", err)
	}

	low, high := 15, 30
	for range employeesNumber {
		if err := processEmployee(&low, &high); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var departmentsNumber int

	if _, err := fmt.Scan(&departmentsNumber); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read number of departments: %v\n", err)

		return
	}

	for range departmentsNumber {
		if err := processDepartment(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}
