package config

type Config struct {
    Environment string `yaml:"environment"`
    LogLevel    string `yaml:"log_level"`
}

func GetConfig() Config {
    var cfg Config
    return cfg
}
