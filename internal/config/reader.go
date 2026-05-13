package config

import (
	"fmt"
	"os"

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

	return &config, nil
}
