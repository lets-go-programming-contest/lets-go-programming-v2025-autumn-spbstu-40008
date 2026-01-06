package main

import (
	"fmt"

	"ivan.manyk/task-8/config"
)

func main() {
    cfg := config.LoadConf()
    if cfg.Environment == "" && cfg.LogLevel == "" {
        fmt.Println("error")
        return
    }
    fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}