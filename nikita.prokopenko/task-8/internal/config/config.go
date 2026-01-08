package config
import (
	"fmt"
	"gopkg.in/yaml.v3"
)
type Settings struct {
	AppStatus   string `yaml:"app_status"`
	ReportLevel string `yaml:"report_level"`
}
func Load() (*Settings, error) {
	var s Settings
	if err := yaml.Unmarshal(yamlData, &s); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return &s, nil
}
