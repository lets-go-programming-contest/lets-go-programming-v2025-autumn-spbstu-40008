package main

import (
	"fmt"

	"github.com/LAffey26/task-8/internal/config"
)

func main() {
	cfg := config.Load()
	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
