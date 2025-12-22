package config

type Config struct {
    Environment string `yaml:"environment"`
    LogLevel    string `yaml:"log_level"`
}

var cfg Config

func GetConfig() Config {
    return cfg
}
