package main

import (
	"errors"
	"fmt"
)

const (
	InitialMinTemp = 15
	InitialMaxTemp = 30
)

var ErrInvalidOperator = errors.New("invalid operator")

type TemperatureManager struct {
	minTemp int
	maxTemp int
}

func NewTemperatureManager() *TemperatureManager {
	return &TemperatureManager{
		minTemp: InitialMinTemp,
		maxTemp: InitialMaxTemp,
	}
}

func (tm *TemperatureManager) Update(operator string, temp int) error {
	switch operator {
	case ">=":
		if temp > tm.minTemp {
			tm.minTemp = temp
		}
	case "<=":
		if temp < tm.maxTemp {
			tm.maxTemp = temp
		}
	default:
		return ErrInvalidOperator
	}

	return nil
}

func (tm *TemperatureManager) GetOptimalTemp() int {
	if tm.minTemp > tm.maxTemp {
		return -1
	}

	return tm.minTemp
}

func processDepartment() error {
	var workers int
	if _, err := fmt.Scan(&workers); err != nil {
		return fmt.Errorf("error when reading number of workers: %w", err)
	}

	manager := NewTemperatureManager()

	for range workers {
		var (
			operator string
			temp     int
		)

		if _, err := fmt.Scan(&operator, &temp); err != nil {
			return fmt.Errorf("error when reading worker input: %w", err)
		}

		if err := manager.Update(operator, temp); err != nil {
			return fmt.Errorf("error when updating temperature manager: %w", err)
		}

		result := manager.GetOptimalTemp()
		fmt.Println(result)
	}

	return nil
}

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Printf("error when reading number of departments: %v\n", err)

		return
	}

	for range departments {
		if err := processDepartment(); err != nil {
			fmt.Printf("error when processing department: %v\n", err)

			return
		}
	}
}
