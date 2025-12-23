package main

import (
	"fmt"

	"esko.dana/task-8/config"
)

func main() {
	conf := config.Load()

	fmt.Print(conf.Environment + " " + conf.LogLevel)
}
