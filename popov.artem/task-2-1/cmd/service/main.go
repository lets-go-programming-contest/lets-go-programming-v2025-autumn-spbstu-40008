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

// Возвращаем результат и ошибку
func updateRange(input string, minTemp, maxTemp *int) (int, error) {
	const minLen = 2
	if len(input) < minLen {
		return 0, nil
	}

	prefix := input[:2]
	if strings.HasPrefix(prefix, ">=") {
		numStr := strings.TrimSpace(input[2:])
		if numStr != "" && numStr[0] == ' ' {
			numStr = strings.TrimSpace(numStr[1:])
		}
		value, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, fmt.Errorf("ошибка при парсинге числа в >=: %w", err)
		}
		if value > *minTemp {
			*minTemp = value
		}
	} else if strings.HasPrefix(prefix, "<=") {
		numStr := strings.TrimSpace(input[2:])
		if numStr != "" && numStr[0] == ' ' {
			numStr = strings.TrimSpace(numStr[1:])
		}
		value, err := strconv.Atoi(numStr)
		if err != nil {
			return 0, fmt.Errorf("ошибка при парсинге числа в <=: %w", err)
		}
		if value < *maxTemp {
			*maxTemp = value
		}
	}

	if *minTemp <= *maxTemp {
		return *minTemp, nil
	}
	return -1, nil
}

func main() {
	scanner := bufio.NewReader(os.Stdin)

	testCount, err := readInt(scanner)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < testCount; i++ {
		conditionCount, err := readInt(scanner)
		if err != nil {
			log.Fatal(err)
		}

		min, max := 15, 30

		for j := 0; j < conditionCount; j++ {
			line, err := readLine(scanner)
			if err != nil {
				log.Fatal(err)
			}

			result, err := updateRange(line, &min, &max)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
		}
	}
}