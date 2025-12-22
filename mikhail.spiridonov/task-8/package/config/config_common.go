package config

type Config struct {
    Environment string
    LogLevel    string
}

var cfg Config

func GetConfig() Config {
    return cfg
}
