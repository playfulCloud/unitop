package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const DefaultConfigYAML = `mode: all # selected | all
refresh_interval: 5s

services:
  - docker.service
  - bluetooth.service
  - NetworkManager.service
  - ssh.service
  - dbus.service
  - systemd-journald.service
  - systemd-logind.service
  - cron.service
  - cups.service
  - avahi-daemon.service
  - firewalld.service
  - polkit.service
  - udisks2.service
  - ModemManager.service
  - wpa_supplicant.service
  - systemd-resolved.service

discovery:
  include:
    - "*.service"
  exclude:
    - "systemd-*"
    - "user@*.service"
    - "getty@*.service"
    - "autovt@*.service"
  states:
    - disabled
    - enabled
    - enabled-runtime
    - linked
    - linked-runtime
`

func DefaultConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user config directory: %w", err)
	}

	return filepath.Join(configDir, "unitop", "unitop.yaml"), nil
}

func EnsureConfigFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to inspect config %q: %w", path, err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory for %q: %w", path, err)
	}

	if err := os.WriteFile(path, []byte(DefaultConfigYAML), 0644); err != nil {
		return fmt.Errorf("failed to create default config %q: %w", path, err)
	}

	return nil
}

func ReadOrCreateConfig(path string) (*AppConfig, error) {
	if err := EnsureConfigFile(path); err != nil {
		return nil, err
	}

	return ReadConfig(path)
}
