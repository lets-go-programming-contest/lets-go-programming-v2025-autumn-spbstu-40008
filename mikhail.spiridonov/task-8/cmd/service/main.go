package main

import (
	"fmt"
	"os"
	"strings"
	"config"
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

func main() {
	cfg := config.GetConfig()

	if err := validateConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(cfg.Environment) == "" {
		fmt.Fprintln(os.Stderr, "Error: environment is empty")
		os.Exit(1)
	}

	if strings.TrimSpace(cfg.LogLevel) == "" {
		fmt.Fprintln(os.Stderr, "Error: log_level is empty")
		os.Exit(1)
	}

	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)

	if cfg.Environment == "dev" && cfg.LogLevel != "debug" {
		fmt.Fprintf(os.Stderr, "Warning: dev environment usually uses debug log level, got %s\n", cfg.LogLevel)
	}

	if cfg.Environment == "prod" && cfg.LogLevel == "debug" {
		fmt.Fprintln(os.Stderr, "Warning: prod environment should not use debug log level for security reasons")
	}
}
