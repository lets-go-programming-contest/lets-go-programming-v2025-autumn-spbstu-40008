//nolint:gofumpt
package main

import (
	"fmt"

	"task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%%s %%s\\n", cfg.Environment, cfg.LogLevel)
}
