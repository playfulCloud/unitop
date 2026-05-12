package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/playfulCloud/unitop/internal/systemd"
)

type tickMsg time.Time

type Model struct {
	collector *systemd.Collector
	table     table.Model
	err       error
	interval  time.Duration
}

func NewModel(collector *systemd.Collector, interval time.Duration) Model {
	columns := []table.Column{
		{Title: "Service", Width: 32},
		{Title: "Load", Width: 12},
		{Title: "Active", Width: 12},
		{Title: "Sub", Width: 14},
		{Title: "PID", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	styles := table.DefaultStyles()

	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)

	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	t.SetStyles(styles)

	return Model{
		collector: collector,
		table:     t,
		interval:  interval,
	}
}

func (m Model) Init() tea.Cmd {
	return tick(m.interval)
}

func tick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		err := m.collector.MonitorState()
		if err != nil {
			m.err = err
		} else {
			m.err = nil
		}

		m.table.SetRows(m.buildRows())

		return m, tick(m.interval)

	case tea.KeyMsg:
		key := msg

		switch key.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) buildRows() []table.Row {
	entries := m.collector.Store.GetServiceEntries()

	rows := make([]table.Row, 0, len(entries))

	for serviceName, entry := range entries {
		rows = append(rows, table.Row{
			serviceName,
			entry.Params["LoadState"],
			entry.Params["ActiveState"],
			entry.Params["SubState"],
			entry.Params["MainPID"],
		})
	}

	return rows
}

func (m Model) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1).
		Render("⚡ Unitop - systemd service monitor")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(m.table.View())

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("↑/↓ navigate • q quit • auto-refresh enabled")

	errorText := ""
	if m.err != nil {
		errorText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Render(fmt.Sprintf("\nError: %v\n", m.err))
	}

	return title + "\n\n" + box + errorText + "\n" + footer + "\n"
}
