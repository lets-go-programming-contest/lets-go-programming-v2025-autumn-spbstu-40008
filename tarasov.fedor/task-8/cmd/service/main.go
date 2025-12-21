package main

import (
	"fmt"
	"github.com/task-8/config"
)

func main() {
	cfg := config.GetConfig()

	fmt.Println(cfg.Environment, cfg.LogLevel)
}
