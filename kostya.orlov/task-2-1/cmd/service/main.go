package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var n int
	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		var k int
		fmt.Scan(&k)

		lower := 15
		upper := 30

		results := make([]int, 0, k)

		for j := 0; j < k; j++ {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)

			if line == "" {
				j--
				continue
			}

			parts := strings.Fields(line)
			if len(parts) != 2 {
				fmt.Println("Ошибка формата ввода")
				return
			}

			op := parts[0]
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Ошибка числа")
				return
			}

			if op == "<=" {
				if val < upper {
					upper = val
				}
			} else if op == ">=" {
				if val > lower {
					lower = val
				}
			}

			if lower > upper {
				results = append(results, -1)
			} else {
				results = append(results, lower)
			}
		}

		for _, r := range results {
			fmt.Println(r)
		}
	}
}
