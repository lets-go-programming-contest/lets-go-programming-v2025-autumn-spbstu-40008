package config

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load() (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
