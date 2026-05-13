package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfig(t *testing.T) {
	validConfig := `
mode: all
refresh_interval: 5s

services:
  - docker.service

properties:
  - ActiveState
  - SubState

discovery:
  include:
    - "*.service"
  exclude:
    - "systemd-*"
  states:
    - enabled
    - disabled
`

	invalidConfig := `
refresh_interval: 5s

services:
  - docker.service
  - [broken

properties:
  - ActiveState
  - SubState
`

	invalidDiscoveryPatternConfig := `
refresh_interval: 5s

services:
  - docker.service

properties:
  - ActiveState

discovery:
  include:
    - "[broken"
`

	invalidModeConfig := `
mode: broken
refresh_interval: 5s

services:
  - docker.service
`

	emptySelectedServicesConfig := `
mode: selected
refresh_interval: 5s

services: []
`

	tests := []struct {
		name         string
		content      string
		path         string
		wantErr      bool
		wantRefresh  string
		wantServices int
		wantMode     string
		wantStates   int
	}{
		{
			name:         "successful parse",
			content:      validConfig,
			path:         "unitop.yaml",
			wantErr:      false,
			wantRefresh:  "5s",
			wantServices: 1,
			wantMode:     "all",
			wantStates:   2,
		},
		{
			name:    "wrong path",
			path:    "wrong-path.yaml",
			wantErr: true,
		},
		{
			name:    "invalid yaml",
			content: invalidConfig,
			path:    "unitop.yaml",
			wantErr: true,
		},
		{
			name:    "invalid discovery pattern",
			content: invalidDiscoveryPatternConfig,
			path:    "unitop.yaml",
			wantErr: true,
		},
		{
			name:    "invalid mode",
			content: invalidModeConfig,
			path:    "unitop.yaml",
			wantErr: true,
		},
		{
			name:    "selected mode without services",
			content: emptySelectedServicesConfig,
			path:    "unitop.yaml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var path string

			if tt.content != "" {
				dir := t.TempDir()
				path = filepath.Join(dir, tt.path)

				err := os.WriteFile(path, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("failed to write test config: %v", err)
				}
			} else {
				path = tt.path
			}

			cfg, err := ReadConfig(path)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if cfg != nil {
					t.Fatalf("expected config to be nil, got %+v", cfg)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if cfg == nil {
				t.Fatal("expected config, got nil")
			}

			if cfg.RefreshInterval != tt.wantRefresh {
				t.Errorf("expected refresh interval %s, got %s", tt.wantRefresh, cfg.RefreshInterval)
			}

			if len(cfg.ServiceNames) != tt.wantServices {
				t.Errorf("expected %d services, got %d", tt.wantServices, len(cfg.ServiceNames))
			}

			if cfg.Mode != tt.wantMode {
				t.Errorf("expected mode %s, got %s", tt.wantMode, cfg.Mode)
			}

			if len(cfg.Discovery.States) != tt.wantStates {
				t.Errorf("expected %d discovery states, got %d", tt.wantStates, len(cfg.Discovery.States))
			}
		})
	}
}
