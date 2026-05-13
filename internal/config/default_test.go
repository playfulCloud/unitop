package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureConfigFileCreatesDefaultConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unitop", "unitop.yaml")

	if err := EnsureConfigFile(path); err != nil {
		t.Fatalf("expected default config to be created, got %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected created config to be readable, got %v", err)
	}

	if string(content) != DefaultConfigYAML {
		t.Fatalf("expected default config content, got %q", string(content))
	}
}

func TestEnsureConfigFileDoesNotOverwriteExistingConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unitop.yaml")
	existingConfig := []byte("mode: selected\nservices:\n  - docker.service\n")

	if err := os.WriteFile(path, existingConfig, 0644); err != nil {
		t.Fatalf("failed to write existing config: %v", err)
	}

	if err := EnsureConfigFile(path); err != nil {
		t.Fatalf("expected existing config to be accepted, got %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected config to be readable, got %v", err)
	}

	if string(content) != string(existingConfig) {
		t.Fatalf("expected existing config to be preserved, got %q", string(content))
	}
}

func TestReadOrCreateConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unitop.yaml")

	cfg, err := ReadOrCreateConfig(path)
	if err != nil {
		t.Fatalf("expected config to be created and read, got %v", err)
	}

	if cfg.Mode != "all" {
		t.Fatalf("expected default mode all, got %q", cfg.Mode)
	}

	if cfg.RefreshInterval != "5s" {
		t.Fatalf("expected default refresh interval 5s, got %q", cfg.RefreshInterval)
	}
}
