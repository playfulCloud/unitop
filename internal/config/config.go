package config

type AppConfig struct {
	Mode            string          `yaml:"mode"`
	Properties      []string        `yaml:"properties"`
	RefreshInterval string          `yaml:"refresh_interval"`
	ServiceNames    []string        `yaml:"services"`
	Discovery       DiscoveryConfig `yaml:"discovery"`
}

type DiscoveryConfig struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
	States  []string `yaml:"states"`
}
