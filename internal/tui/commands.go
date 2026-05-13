package tui

import (
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/playfulCloud/unitop/internal/systemd"
)

func tick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func monitorStateCmd(manager *systemd.SystemdManager) tea.Cmd {
	return func() tea.Msg {
		err := manager.MonitorState()
		return monitorDoneMsg{err: err}
	}
}

func enterJournalctlCmd(serviceID string) tea.Cmd {
	command := systemd.BuildJournalctlCommand(serviceID)
	cmd := exec.Command(command.Name, command.Args...)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return journalctlDoneMsg{err: err}
	})
}

func executeActionCmd(
	manager *systemd.SystemdManager,
	serviceID string,
	action systemd.ServiceAction,
) tea.Cmd {
	return func() tea.Msg {
		err := manager.ExecuteAction(serviceID, action)
		return actionDoneMsg{err: err}
	}
}
