package config

import (
	"strings"
)

type Config struct {
	Environment string
	LogLevel    string
}

const splitParts = 2

func ParseConfig(data []byte) Config {
	cfg := Config{
		Environment: "",
		LogLevel:    "",
	}
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", splitParts)
		if len(parts) != splitParts {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "\"")

		switch key {
		case "environment":
			cfg.Environment = val
		case "log_level":
			cfg.LogLevel = val
		}
	}

	return cfg
}
