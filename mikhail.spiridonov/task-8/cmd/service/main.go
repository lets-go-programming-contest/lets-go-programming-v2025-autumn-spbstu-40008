package main

import (
	"fmt"

	"github.com/mordw1n/task-8/internal/config"
)

func main() {
	cfg := config.GetConfig()
	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
