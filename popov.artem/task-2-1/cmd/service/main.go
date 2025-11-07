package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	_, _ = strconv.Atoi(scanner.Text())

	scanner.Scan()
	k, _ := strconv.Atoi(scanner.Text())

	minReq := 15
	maxReq := 30

	for i := 0; i < k; i++ {
		scanner.Scan()
		line := scanner.Text()
		parts := strings.Split(line, " ")
		sign := parts[0]
		temp, _ := strconv.Atoi(parts[1])

		if sign == ">=" {
			if temp > minReq {
				minReq = temp
			}
		} else if sign == "<=" {
			if temp < maxReq {
				maxReq = temp
			}
		}

		if minReq <= maxReq {
			fmt.Println(minReq)
		} else {
			fmt.Println(-1)
		}
	}
}
