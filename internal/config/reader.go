package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadConfig(path string) (*AppConfig, error) {
	dataFromConfigFile, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read a config: %w", err)
	}

	var config AppConfig

	err = yaml.Unmarshal(dataFromConfigFile, &config)

	if err != nil {
		return nil, fmt.Errorf("failed to parse a config from yaml: %w", err)
	}

	if err := validateDiscoveryPatterns(config.Discovery); err != nil {
		return nil, err
	}

	if err := validateMode(config); err != nil {
		return nil, err
	}

	config.Mode = strings.ToLower(strings.TrimSpace(config.Mode))

	return &config, nil
}

func validateMode(config AppConfig) error {
	mode := strings.ToLower(strings.TrimSpace(config.Mode))
	switch mode {
	case ModeAll:
		return nil
	case ModeSelected:
		if len(config.ServiceNames) == 0 {
			return fmt.Errorf("services must contain at least one service when mode is %q", ModeSelected)
		}

		return nil
	default:
		return fmt.Errorf("invalid mode %q: expected %q or %q", config.Mode, ModeAll, ModeSelected)
	}
}

func validateDiscoveryPatterns(discovery DiscoveryConfig) error {
	for _, pattern := range discovery.Include {
		if _, err := filepath.Match(pattern, ""); err != nil {
			return fmt.Errorf("invalid discovery include pattern %q: %w", pattern, err)
		}
	}

	for _, pattern := range discovery.Exclude {
		if _, err := filepath.Match(pattern, ""); err != nil {
			return fmt.Errorf("invalid discovery exclude pattern %q: %w", pattern, err)
		}
	}

	return nil
}
