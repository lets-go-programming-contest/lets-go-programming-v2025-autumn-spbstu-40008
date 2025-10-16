package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func readLine() (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func main() {
	var n int

	if _, err := fmt.Scan(&n); err != nil {
		os.Exit(0)
	}

	for i := 0; i < n; i++ {
		var k int

		if _, err := fmt.Scan(&k); err != nil {
			os.Exit(0)
		}

		lower := 15
		upper := 30

		for j := 0; j < k; j++ {
			temp, err := readLine()

			if err != nil || temp == "" {
				os.Exit(0)
			}

			parts := strings.Fields(temp)
			if len(parts) != 2 {
				upper = -1
				lower = 0
			} else {
				operator := parts[0]
				degreeStr := parts[1]

				degree, err := strconv.Atoi(degreeStr)
				if err != nil {
					os.Exit(0)
				}

				switch operator {
				case "<=":
					if degree < upper {
						upper = degree
					}
				case ">=":
					if degree > lower {
						lower = degree
					}
				}
			}

			if lower > upper {
				fmt.Println(-1)
			} else {
				fmt.Println(lower)
			}
		}
	}
}
