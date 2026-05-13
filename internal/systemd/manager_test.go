package systemd

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/playfulCloud/unitop/internal/model"
	"github.com/playfulCloud/unitop/internal/store"
)

func TestSystemdManagerMonitorStateSuccess(t *testing.T) {
	serviceStore := store.NewServiceStore(
		[]string{"docker.service"},
		[]string{"ID", "LoadState", "ActiveState"},
	)

	manager := NewSystemdManager(
		serviceStore,
		[]string{"ID", "LoadState", "ActiveState"},
	)

	manager.Execute = func(command model.Command) (string, error) {
		return `
			ID=docker.service
			LoadState=loaded
			ActiveState=active
		`, nil
	}

	err := manager.MonitorState()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	entry, exists := serviceStore.GetServiceEntry("docker.service")
	if !exists {
		t.Fatal("expected docker.service to exist")
	}

	if entry.Params["ID"] != "docker.service" {
		t.Fatalf("expected ID docker.service, got %s", entry.Params["ID"])
	}

	if entry.Params["LoadState"] != "loaded" {
		t.Fatalf("expected LoadState loaded, got %s", entry.Params["LoadState"])
	}

	if entry.Params["ActiveState"] != "active" {
		t.Fatalf("expected ActiveState active, got %s", entry.Params["ActiveState"])
	}
}

func TestSystemdManagerMonitorStateReturnsErrorWhenExecuteFails(t *testing.T) {
	serviceStore := store.NewServiceStore(
		[]string{"docker.service"},
		[]string{"ID", "LoadState", "ActiveState"},
	)

	manager := NewSystemdManager(
		serviceStore,
		[]string{"ID", "LoadState", "ActiveState"},
	)

	manager.Execute = func(command model.Command) (string, error) {
		return "", fmt.Errorf("boom")
	}

	err := manager.MonitorState()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to monitor docker.service") {
		t.Fatalf("expected error to contain service ID, got %v", err)
	}
}

func TestSystemdManagerExecuteActionSuccess(t *testing.T) {
	serviceStore := store.NewServiceStore(
		[]string{"docker.service"},
		[]string{"ID", "LoadState", "ActiveState"},
	)

	manager := NewSystemdManager(
		serviceStore,
		[]string{"ID", "LoadState", "ActiveState"},
	)

	var executedCommand model.Command

	manager.Execute = func(command model.Command) (string, error) {
		executedCommand = command
		return "", nil
	}

	err := manager.ExecuteAction("docker.service", RestartAction)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if executedCommand.Name != "systemctl" {
		t.Fatalf("expected command name systemctl, got %s", executedCommand.Name)
	}

	expectedArgs := []string{"restart", "docker.service"}

	if !reflect.DeepEqual(expectedArgs, executedCommand.Args) {
		t.Fatalf("expected args %v, got %v", expectedArgs, executedCommand.Args)
	}
}
