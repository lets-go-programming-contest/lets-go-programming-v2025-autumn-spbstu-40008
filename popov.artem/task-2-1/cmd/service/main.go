package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInt(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения: %w", err)
	}
	line = strings.TrimSpace(line)
	num, err := strconv.Atoi(line)
	if err != nil {
		return 0, fmt.Errorf("ошибка конвертации: %w", err)
	}
	return num, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("ошибка чтения: %w", err)
	}
	return strings.TrimSpace(line), nil
}

func parseCondition(input string, currentMin, currentMax int) (int, int, error) {
	if len(input) < 2 {
		return currentMin, currentMax, nil
	}

	prefix := input[:2]
	numStr := strings.TrimSpace(input[2:])

	if prefix == ">=" {
		value, err := strconv.Atoi(numStr)
		if err != nil {
			return currentMin, currentMax, fmt.Errorf("ошибка при парсинге числа в >=: %w", err)
		}
		if value > currentMin {
			currentMin = value
		}
	} else if prefix == "<=" {
		value, err := strconv.Atoi(numStr)
		if err != nil {
			return currentMin, currentMax, fmt.Errorf("ошибка при парсинге числа в <=: %w", err)
		}
		if value < currentMax {
			currentMax = value
		}
	}
	return currentMin, currentMax, nil
}

func checkRange(min, max int) int {
	if min <= max {
		return min
	}
	return -1
}

func main() {
	scanner := bufio.NewReader(os.Stdin)

	testCount, err := readInt(scanner)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < testCount; i++ {
		_ = i
		conditionCount, err := readInt(scanner)
		if err != nil {
			log.Fatal(err)
		}

		minVal, maxVal := 15, 30

		for j := 0; j < conditionCount; j++ {
			_ = j
			line, err := readLine(scanner)
			if err != nil {
				log.Fatal(err)
			}

			minVal, maxVal, err = parseCondition(line, minVal, maxVal)
			if err != nil {
				log.Fatal(err)
			}

			result := checkRange(minVal, maxVal)
			fmt.Println(result)
		}
	}
}