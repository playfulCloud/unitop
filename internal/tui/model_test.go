package tui

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	domainmodel "github.com/playfulCloud/unitop/internal/model"
	"github.com/playfulCloud/unitop/internal/store"
	"github.com/playfulCloud/unitop/internal/systemd"
)

func TestTickSchedulesMonitorWithoutRunningSynchronously(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		t.Fatal("expected monitor command not to run during Update")
		return "", nil
	})
	m := NewModel(manager, time.Second)

	updated, cmd := m.Update(tickMsg(time.Now()))
	updatedModel := updated.(Model)

	if !updatedModel.monitoring {
		t.Fatal("expected model to mark monitoring in flight")
	}

	if cmd == nil {
		t.Fatal("expected tick to return a command")
	}

	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected batch command, got %T", msg)
	}

	if len(batch) != 2 {
		t.Fatalf("expected monitor and next tick commands, got %d", len(batch))
	}
}

func TestMonitorDoneClearsMonitoringAndNormalizesSelection(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})
	m := NewModel(manager, time.Second)
	m.monitoring = true

	updated, cmd := m.Update(monitorDoneMsg{})
	updatedModel := updated.(Model)

	if cmd != nil {
		t.Fatalf("expected no command after monitor completes, got %T", cmd)
	}

	if updatedModel.monitoring {
		t.Fatal("expected monitoring flag to be cleared")
	}

	if updatedModel.selectedServiceID != "docker.service" {
		t.Fatalf("expected selection to be normalized, got %q", updatedModel.selectedServiceID)
	}
}

func TestActionDoneDoesNotScheduleTickWhileMonitorIsRunning(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		t.Fatal("expected no command execution")
		return "", nil
	})
	m := NewModel(manager, time.Second)
	m.monitoring = true

	updated, cmd := m.Update(actionDoneMsg{})
	updatedModel := updated.(Model)

	if cmd != nil {
		t.Fatalf("expected no command while monitor is running, got %T", cmd)
	}

	if !updatedModel.monitoring {
		t.Fatal("expected existing monitor to remain in flight")
	}
}

func TestActionDoneWithErrorDoesNotStartMonitor(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		t.Fatal("expected no command execution")
		return "", nil
	})
	m := NewModel(manager, time.Second)
	actionErr := errors.New("action failed")

	updated, cmd := m.Update(actionDoneMsg{err: actionErr})
	updatedModel := updated.(Model)

	if cmd != nil {
		t.Fatalf("expected no command after failed action, got %T", cmd)
	}

	if !errors.Is(updatedModel.err, actionErr) {
		t.Fatalf("expected action error to be stored, got %v", updatedModel.err)
	}
}

func TestActionDoneSchedulesSingleImmediateMonitor(t *testing.T) {
	calls := 0
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		calls++
		return "LoadState=loaded\nActiveState=active\n", nil
	})
	m := NewModel(manager, time.Second)

	updated, cmd := m.Update(actionDoneMsg{})
	updatedModel := updated.(Model)

	if !updatedModel.monitoring {
		t.Fatal("expected action completion to start an immediate monitor")
	}

	if cmd == nil {
		t.Fatal("expected action completion to return monitor command")
	}

	msg := cmd()
	if _, ok := msg.(monitorDoneMsg); !ok {
		t.Fatalf("expected monitor done message, got %T", msg)
	}

	if calls != 1 {
		t.Fatalf("expected one monitor execution, got %d", calls)
	}
}

func newTestManager(
	t *testing.T,
	execute func(command domainmodel.Command) (string, error),
) *systemd.SystemdManager {
	t.Helper()

	serviceStore := store.NewServiceStore(
		[]string{"docker.service"},
		[]string{"LoadState", "ActiveState"},
	)

	manager := systemd.NewSystemdManager(
		serviceStore,
		[]string{"LoadState", "ActiveState"},
	)
	manager.Execute = execute

	return manager
}
