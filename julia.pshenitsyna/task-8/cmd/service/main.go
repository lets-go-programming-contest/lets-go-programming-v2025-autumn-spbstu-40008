package main

import (
	"fmt"
	"github.com/julia.pshenitsyna/task-8/internal/config"
	"log"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	fmt.Printf("%s %s\n", conf.Environment, conf.LogLevel)
}
