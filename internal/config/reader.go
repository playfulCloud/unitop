package config

import (
	"fmt"
	"os"
	"path/filepath"

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

	return &config, nil
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
