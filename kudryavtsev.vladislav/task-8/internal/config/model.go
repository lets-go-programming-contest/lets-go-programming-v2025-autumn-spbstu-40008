package config

type Settings struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}
