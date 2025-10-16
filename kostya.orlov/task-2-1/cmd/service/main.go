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

		for j := 0; j < k; j++ {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)

			if line == "" {
				j--
				continue
			}

			parts := strings.Fields(line)
			if len(parts) != 2 {
				fmt.Println("Error input")
				return
			}

			op := parts[0]
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error num")
				return
			}

			switch op {
			case "<=":
				upper = min(upper, val)
			case ">=":
				lower = max(lower, val)
			default:
				continue
			}

			if lower > upper {
				fmt.Println(-1)
			} else {
				fmt.Print(lower)
			}
		}
	}
}
