package config

const (
	ModeAll      = "all"
	ModeSelected = "selected"
)

type AppConfig struct {
	Mode            string          `yaml:"mode"`
	RefreshInterval string          `yaml:"refresh_interval"`
	ServiceNames    []string        `yaml:"services"`
	Discovery       DiscoveryConfig `yaml:"discovery"`
}

type DiscoveryConfig struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
	States  []string `yaml:"states"`
}
