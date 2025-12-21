package main

import(
	"fmt"
	"github.com/julia.pshenitsyna/task-8/internal/config"
)

func main() {
	configg := config.GetConfig()
	fmt.Printf("%s %s\n", configg.Environment, configg.LogLevel)
}