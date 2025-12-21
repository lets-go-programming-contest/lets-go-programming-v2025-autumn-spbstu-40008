package main

import (
	"fmt"
	"log"

	"github.com/julia.pshenitsyna/task-8/internal/config"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	fmt.Print(conf.Environment, " ",conf.LogLevel)
}
