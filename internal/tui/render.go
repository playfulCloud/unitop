package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
)

func (m Model) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1).
		Render("⚡ Unitop - systemd service monitor")

	box := lipgloss.NewStyle().
		Render(m.renderTable())

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(m.renderFooter())

	errorText := ""
	if m.err != nil {
		errorText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Render(fmt.Sprintf("\nError: %v\n", m.err))
	}

	return title + "\n\n" + box + errorText + "\n" + footer + "\n"
}

func (m Model) renderTable() string {
	entries := m.systemdManager.Store.GetServiceEntries()
	serviceNames := m.sortedServiceNames()

	if len(serviceNames) == 0 {
		if strings.TrimSpace(m.filterText) != "" {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("244")).
				Italic(true).
				Render(fmt.Sprintf("No services matching /%s", m.filterText))
		}

		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true).
			Render("No services to display")
	}

	start := m.viewportOffset
	end := min(start+m.tableHeight, len(serviceNames))

	rows := make([][]string, 0, end-start)

	for _, serviceName := range serviceNames[start:end] {
		entry := entries[serviceName]

		rows = append(rows, []string{
			serviceName,
			emptyAsDash(entry.Params["LoadState"]),
			emptyAsDash(entry.Params["ActiveState"]),
			emptyAsDash(entry.Params["SubState"]),
			formatPID(entry.Params["MainPID"]),
		})
	}

	return ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		Headers("Service", "Load", "Active", "Sub", "PID").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			base := lipgloss.NewStyle().Padding(0, 1)

			if row == ltable.HeaderRow {
				return base.
					Bold(true).
					Foreground(lipgloss.Color("252"))
			}

			serviceID := rows[row][0]
			selected := serviceID == m.selectedServiceID

			if selected {
				return base.
					Background(lipgloss.Color("57")).
					Foreground(lipgloss.Color("229")).
					Bold(true)
			}

			if col == 2 {
				activeState := rows[row][2]
				return styleForActiveState(base, activeState)
			}

			if col == 4 {
				pid := rows[row][4]
				if pid == "-" {
					return base.Foreground(lipgloss.Color("244"))
				}

				return base.Foreground(lipgloss.Color("86"))
			}

			return base
		}).
		String()
}

func (m Model) renderFooter() string {
	if m.actionInFlight {
		return fmt.Sprintf(
			"Running: %s %s",
			m.pendingAction,
			m.pendingServiceID,
		)
	}

	if m.filterMode {
		return fmt.Sprintf(
			"Filter: /%s | enter: apply | esc: close filter | backspace: delete",
			m.filterText,
		)
	}

	if strings.TrimSpace(m.filterText) != "" {
		return fmt.Sprintf(
			"Filter: /%s | /: new filter | esc: clear filter | %s",
			m.filterText,
			footerText(),
		)
	}

	return footerText()
}
