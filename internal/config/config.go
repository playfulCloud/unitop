package config

type AppConfig struct {
	Mode            string   `yaml:"mode"`
	Properties      []string `yaml:"properties"`
	RefreshInterval string   `yaml:"refresh_interval"`
	ServiceNames    []string `yaml:"services"`
}
