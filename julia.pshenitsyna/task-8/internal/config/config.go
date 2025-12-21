package config

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

var appConfig *Config

func GetConfig() *Config {
	if appConfig == nil {
		panic("config not initialized")
	}
	return appConfig
}