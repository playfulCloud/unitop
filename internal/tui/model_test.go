package tui

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	domainmodel "github.com/playfulCloud/unitop/internal/model"
	"github.com/playfulCloud/unitop/internal/store"
	"github.com/playfulCloud/unitop/internal/systemd"
)

func TestInitSchedulesInitialMonitorAndTick(t *testing.T) {
	var calls atomic.Int32
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		calls.Add(1)
		return "LoadState=loaded\nActiveState=active\n", nil
	})
	m := NewModel(manager, time.Second)

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected init command")
	}

	msg := cmd()
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected batch command, got %T", msg)
	}

	if len(batch) != 2 {
		t.Fatalf("expected monitor and tick commands, got %d", len(batch))
	}

	monitorMsg := batch[0]()
	if _, ok := monitorMsg.(monitorDoneMsg); !ok {
		t.Fatalf("expected monitor done message, got %T", monitorMsg)
	}

	if calls.Load() != 1 {
		t.Fatalf("expected one monitor execution, got %d", calls.Load())
	}
}

func TestTickSchedulesMonitorWithoutRunningSynchronously(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		t.Fatal("expected monitor command not to run during Update")
		return "", nil
	})
	m := NewModel(manager, time.Second)
	m.monitoring = false

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

	if updatedModel.actionInFlight {
		t.Fatal("expected action state to be cleared after failed action")
	}
}

func TestActionDoneSchedulesSingleImmediateMonitor(t *testing.T) {
	calls := 0
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		calls++
		return "LoadState=loaded\nActiveState=active\n", nil
	})
	m := NewModel(manager, time.Second)
	m.monitoring = false

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

	if updatedModel.actionInFlight {
		t.Fatal("expected action state to be cleared after successful action")
	}
}

func TestActionKeyMarksActionInFlight(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})
	m := NewModel(manager, time.Second)
	m.selectedServiceID = "docker.service"

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	updatedModel := updated.(Model)

	if cmd == nil {
		t.Fatal("expected action command")
	}

	if !updatedModel.actionInFlight {
		t.Fatal("expected action to be marked in flight")
	}

	if updatedModel.pendingAction != systemd.StopAction {
		t.Fatalf("expected pending action stop, got %q", updatedModel.pendingAction)
	}

	if updatedModel.pendingServiceID != "docker.service" {
		t.Fatalf("expected pending service docker.service, got %q", updatedModel.pendingServiceID)
	}
}

func TestActionKeyIgnoredWhileActionInFlight(t *testing.T) {
	executed := false
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		executed = true
		return "", nil
	})
	m := NewModel(manager, time.Second)
	m.selectedServiceID = "docker.service"
	m.actionInFlight = true
	m.pendingAction = systemd.StopAction
	m.pendingServiceID = "docker.service"

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	updatedModel := updated.(Model)

	if cmd != nil {
		t.Fatalf("expected no command while action is in flight, got %T", cmd)
	}

	if executed {
		t.Fatal("expected no action execution while action is in flight")
	}

	if !updatedModel.actionInFlight {
		t.Fatal("expected action to remain in flight")
	}

	if updatedModel.pendingAction != systemd.StopAction {
		t.Fatalf("expected pending action to remain stop, got %q", updatedModel.pendingAction)
	}
}

func TestActionFooterShowsInFlightAction(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})
	m := NewModel(manager, time.Second)
	m.actionInFlight = true
	m.pendingAction = systemd.RestartAction
	m.pendingServiceID = "docker.service"

	footer := m.renderFooter()

	if footer != "Running: restart docker.service" {
		t.Fatalf("expected running action footer, got %q", footer)
	}
}

func TestSlashEntersFilterModeAndClearsPreviousFilter(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})

	m := NewModel(manager, time.Second)
	m.filterText = "old"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	updatedModel := updated.(Model)

	if !updatedModel.filterMode {
		t.Fatal("expected filter mode to be enabled")
	}

	if updatedModel.filterText != "" {
		t.Fatalf("expected filter text to be cleared, got %q", updatedModel.filterText)
	}
}

func TestFilterModeAppendsTypedCharacters(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})

	m := NewModel(manager, time.Second)
	m.filterMode = true

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	updatedModel := updated.(Model)

	if updatedModel.filterText != "d" {
		t.Fatalf("expected filter text to be updated, got %q", updatedModel.filterText)
	}
}

func TestFilterModeDoesNotExecuteActions(t *testing.T) {
	executed := false

	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		executed = true
		return "", nil
	})

	m := NewModel(manager, time.Second)
	m.filterMode = true
	m.selectedServiceID = "docker.service"

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updatedModel := updated.(Model)

	if cmd != nil {
		t.Fatalf("expected no action command in filter mode, got %T", cmd)
	}

	if executed {
		t.Fatal("expected action not to be executed in filter mode")
	}

	if updatedModel.filterText != "r" {
		t.Fatalf("expected r to be added to filter text, got %q", updatedModel.filterText)
	}
}

func TestFilterEnterLeavesFilterModeButKeepsFilter(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})

	m := NewModel(manager, time.Second)
	m.filterMode = true
	m.filterText = "dock"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updatedModel := updated.(Model)

	if updatedModel.filterMode {
		t.Fatal("expected filter mode to be disabled")
	}

	if updatedModel.filterText != "dock" {
		t.Fatalf("expected filter text to be kept, got %q", updatedModel.filterText)
	}
}

func TestFilterTextFiltersSortedServiceNames(t *testing.T) {
	manager := newTestManagerWithServices(t,
		[]string{"docker.service", "nginx.service", "postgresql.service"},
		func(command domainmodel.Command) (string, error) {
			return "", nil
		},
	)

	m := NewModel(manager, time.Second)
	m.filterText = "gin"

	names := m.sortedServiceNames()

	if len(names) != 1 {
		t.Fatalf("expected one filtered service, got %d", len(names))
	}

	if names[0] != "nginx.service" {
		t.Fatalf("expected nginx.service, got %q", names[0])
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

func newTestManagerWithServices(
	t *testing.T,
	services []string,
	execute func(command domainmodel.Command) (string, error),
) *systemd.SystemdManager {
	t.Helper()

	serviceStore := store.NewServiceStore(
		services,
		[]string{"LoadState", "ActiveState"},
	)

	manager := systemd.NewSystemdManager(
		serviceStore,
		[]string{"LoadState", "ActiveState"},
	)
	manager.Execute = execute

	return manager
}

func TestJournalctlDoesNothingWhenNoServiceSelected(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		t.Fatal("expected no command execution")
		return "", nil
	})

	m := NewModel(manager, time.Second)
	m.selectedServiceID = ""

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	if cmd != nil {
		t.Fatalf("expected no command when no service is selected, got %T", cmd)
	}
}

func TestJournalctlDoneStoresError(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})

	m := NewModel(manager, time.Second)
	journalErr := errors.New("journalctl failed")

	updated, cmd := m.Update(journalctlDoneMsg{err: journalErr})
	updatedModel := updated.(Model)

	if cmd != nil {
		t.Fatalf("expected no command after journalctl completes, got %T", cmd)
	}

	if !errors.Is(updatedModel.err, journalErr) {
		t.Fatalf("expected journalctl error to be stored, got %v", updatedModel.err)
	}
}

func TestFilterBackspaceRemovesLastCharacter(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})

	m := NewModel(manager, time.Second)
	m.filterMode = true
	m.filterText = "dock"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	updatedModel := updated.(Model)

	if updatedModel.filterText != "doc" {
		t.Fatalf("expected filter text to be shortened, got %q", updatedModel.filterText)
	}
}

func TestEscapeClearsAppliedFilter(t *testing.T) {
	manager := newTestManager(t, func(command domainmodel.Command) (string, error) {
		return "", nil
	})

	m := NewModel(manager, time.Second)
	m.filterText = "dock"

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updatedModel := updated.(Model)

	if updatedModel.filterText != "" {
		t.Fatalf("expected filter text to be cleared, got %q", updatedModel.filterText)
	}
}
