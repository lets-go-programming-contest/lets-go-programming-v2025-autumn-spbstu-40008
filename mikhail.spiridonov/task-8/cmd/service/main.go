package main

import (
	"fmt"
	"strings"

	"github.com/task-8/internal/config"
)

func validateConfig(cfg config.Config) error {
	validEnvironments := map[string]bool{
		"dev":  true,
		"prod": true,
	}
	if !validEnvironments[cfg.Environment] {
		return fmt.Errorf("invalid environment: %s", cfg.Environment)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[cfg.LogLevel] {
		return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	return nil
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := config.GetConfig()

	if err := validateConfig(cfg); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		panic("Configuration validation failed")
	}

	if strings.TrimSpace(cfg.Environment) == "" {
		fmt.Println("Error: environment is empty")
		panic("Empty environment")
	}

	if strings.TrimSpace(cfg.LogLevel) == "" {
		fmt.Println("Error: log_level is empty")
		panic("Empty log_level")
	}

	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)

	if cfg.Environment == "dev" && cfg.LogLevel != "debug" {
		fmt.Printf("Warning: dev environment usually uses debug log level, got %s\n", cfg.LogLevel)
	}

	if cfg.Environment == "prod" && cfg.LogLevel == "debug" {
		fmt.Println("Warning: prod environment should not use debug log level for security reasons")
	}
}
