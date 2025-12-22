package config

type Config struct {
    Environment string `yaml:"environment"`
    LogLevel    string `yaml:"log_level"`
}

func Load(data []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal YAML config: %w", err)
	}

	return cfg, nil
}