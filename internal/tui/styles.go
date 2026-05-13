package tui

import (
	"slices"

	"github.com/charmbracelet/lipgloss"
)

func styleForActiveState(base lipgloss.Style, state string) lipgloss.Style {
	switch state {
	case "active":
		return base.Foreground(lipgloss.Color("42")).Bold(true)
	case "failed":
		return base.Foreground(lipgloss.Color("196")).Bold(true)
	case "inactive":
		return base.Foreground(lipgloss.Color("244"))
	case "activating", "deactivating":
		return base.Foreground(lipgloss.Color("214")).Bold(true)
	default:
		return base.Foreground(lipgloss.Color("245"))
	}
}

func contains(values []string, target string) bool {
	return slices.Contains(values, target)
}

func formatPID(pid string) string {
	if pid == "" || pid == "0" {
		return "-"
	}

	return pid
}

func emptyAsDash(value string) string {
	if value == "" {
		return "-"
	}

	return value
}
