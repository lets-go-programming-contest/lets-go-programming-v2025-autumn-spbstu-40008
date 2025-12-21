package config

import (
	"strings"
)

type Config struct {
	Environment string
	LogLevel    string
}

func ParseConfig(data []byte) Config {
	cfg := Config{}
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
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
