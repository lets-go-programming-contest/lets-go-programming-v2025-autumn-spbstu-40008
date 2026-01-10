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
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	for range n {
		var k int
		if _, err := fmt.Fscan(in, &k); err != nil {
			break
		}

		minTemp := initialMinTemp
		maxTemp := initialMaxTemp

		for range k {
			var op string
			var val int

			if _, err := fmt.Fscan(in, &op, &val); err != nil {
				break
			}

			if op == ">=" {
				minTemp = max(minTemp, val)
			} else {
				maxTemp = min(maxTemp, val)
			}

			if minTemp <= maxTemp {
				fmt.Fprintln(out, minTemp)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}
