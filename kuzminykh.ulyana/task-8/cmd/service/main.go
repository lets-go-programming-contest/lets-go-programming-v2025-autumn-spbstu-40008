package main

import (
	"fmt"

	"github.com/kuzminykh.ulyana/task-8/internal/config"
)

func main() {
	cfg := config.Get()

	fmt.Print(cfg.Environment + " " + cfg.LogLevel)
}
