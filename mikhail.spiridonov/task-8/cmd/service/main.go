package main

import (
	"fmt"
	"task-8/package/config"
)

func main() {
	cfg := config.GetConfig()
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
