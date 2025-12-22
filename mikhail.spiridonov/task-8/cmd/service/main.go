package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mordw1n/task-8/internal/config"
)

var (
	errInvalidEnvironment = errors.New("invalid environment")
	errInvalidLogLevel    = errors.New("invalid log level")
)

func validateConfig(cfg config.Config) error {
	validEnvironments := map[string]bool{
		"dev":  true,
		"prod": true,
	}
	if !validEnvironments[cfg.Environment] {
		return fmt.Errorf("%w: %s", errInvalidEnvironment, cfg.Environment)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[cfg.LogLevel] {
		return fmt.Errorf("%w: %s", errInvalidLogLevel, cfg.LogLevel)
	}

	return nil
}

func main() {
	cfg := config.GetConfig()

	if err := validateConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		panic("Configuration validation failed")
	}

	if strings.TrimSpace(cfg.Environment) == "" {
		fmt.Fprintln(os.Stderr, "Error: environment is empty")
		panic("Empty environment")
	}

	if strings.TrimSpace(cfg.LogLevel) == "" {
		fmt.Fprintln(os.Stderr, "Error: log_level is empty")
		panic("Empty log_level")
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)

	if cfg.Environment == "dev" && cfg.LogLevel != "debug" {
		fmt.Fprintf(os.Stderr, "Warning: dev environment usually uses debug log level, got %s\n", cfg.LogLevel)
	}

	if cfg.Environment == "prod" && cfg.LogLevel == "debug" {
		fmt.Fprintln(os.Stderr, "Warning: prod environment should not use debug log level for security reasons")
	}
}
