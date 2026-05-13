package tui

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"

	"github.com/playfulCloud/unitop/internal/systemd"
)

type tickMsg time.Time

type monitorDoneMsg struct {
	err error
}

type actionDoneMsg struct {
	err error
}

type journalctlDoneMsg struct {
	err error
}

type Model struct {
	systemdManager    *systemd.SystemdManager
	err               error
	interval          time.Duration
	selectedServiceID string
	viewportOffset    int
	tableHeight       int
	monitoring        bool
	filterMode        bool
	filterText        string
}

func NewModel(systemdManager *systemd.SystemdManager, interval time.Duration) Model {
	return Model{
		systemdManager: systemdManager,
		interval:       interval,
		tableHeight:    20,
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.monitoring {
			return m, tick(m.interval)
		}

		m.monitoring = true
		return m, tea.Batch(
			monitorStateCmd(m.systemdManager),
			tick(m.interval),
		)

	case monitorDoneMsg:
		m.monitoring = false
		m.err = msg.err
		m.normalizeSelection()
		return m, nil

	case actionDoneMsg:
		m.err = msg.err
		if msg.err != nil {
			return m, nil
		}

		if m.monitoring {
			return m, nil
		}

		m.monitoring = true
		return m, monitorStateCmd(m.systemdManager)

	case journalctlDoneMsg:
		m.err = msg.err
		return m, nil

	case tea.KeyMsg:
		if m.filterMode {
			return m.updateFilter(msg)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			m.moveUp()

		case "down", "j":
			m.moveDown()

		case "/":
			m.filterMode = true
			m.filterText = ""
			m.normalizeSelection()
			return m, nil

		case "esc":
			m.filterText = ""
			m.normalizeSelection()
			return m, nil

		case "l":
			if m.selectedServiceID == "" {
				return m, nil
			}

			return m, enterJournalctlCmd(m.selectedServiceID)

		default:
			if action, ok := actionForKey(msg.String()); ok {
				if m.selectedServiceID == "" {
					return m, nil
				}

				return m, executeActionCmd(
					m.systemdManager,
					m.selectedServiceID,
					action,
				)
			}
		}
	}

	return m, nil
}

func (m Model) updateFilter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.filterMode = false
		return m, nil

	case "enter":
		m.filterMode = false
		m.normalizeSelection()
		return m, nil

	case "backspace":
		if len(m.filterText) > 0 {
			m.filterText = m.filterText[:len(m.filterText)-1]
			m.normalizeSelection()
		}
		return m, nil

	case "ctrl+c":
		return m, tea.Quit

	default:
		if len(msg.String()) == 1 {
			m.filterText += msg.String()
			m.normalizeSelection()
		}

		return m, nil
	}
}

func (m *Model) moveUp() {
	serviceNames := m.sortedServiceNames()
	if len(serviceNames) == 0 {
		return
	}

	currentIndex := m.currentSelectedIndex(serviceNames)

	if currentIndex > 0 {
		currentIndex--
	}

	m.selectedServiceID = serviceNames[currentIndex]
	m.adjustViewport(currentIndex)
}

func (m *Model) moveDown() {
	serviceNames := m.sortedServiceNames()
	if len(serviceNames) == 0 {
		return
	}

	currentIndex := m.currentSelectedIndex(serviceNames)

	if currentIndex < len(serviceNames)-1 {
		currentIndex++
	}

	m.selectedServiceID = serviceNames[currentIndex]
	m.adjustViewport(currentIndex)
}

func (m *Model) normalizeSelection() {
	serviceNames := m.sortedServiceNames()

	if len(serviceNames) == 0 {
		m.selectedServiceID = ""
		m.viewportOffset = 0
		return
	}

	if m.selectedServiceID == "" || !contains(serviceNames, m.selectedServiceID) {
		m.selectedServiceID = serviceNames[0]
		m.viewportOffset = 0
		return
	}

	currentIndex := m.currentSelectedIndex(serviceNames)
	m.adjustViewport(currentIndex)
}

func (m *Model) adjustViewport(selectedIndex int) {
	if selectedIndex < m.viewportOffset {
		m.viewportOffset = selectedIndex
	}

	if selectedIndex >= m.viewportOffset+m.tableHeight {
		m.viewportOffset = selectedIndex - m.tableHeight + 1
	}

	if m.viewportOffset < 0 {
		m.viewportOffset = 0
	}
}

func (m Model) currentSelectedIndex(serviceNames []string) int {
	for index, serviceName := range serviceNames {
		if serviceName == m.selectedServiceID {
			return index
		}
	}

	return 0
}

func (m Model) sortedServiceNames() []string {
	entries := m.systemdManager.Store.GetServiceEntries()

	serviceNames := make([]string, 0, len(entries))
	filter := strings.ToLower(strings.TrimSpace(m.filterText))

	for serviceName := range entries {
		if filter != "" && !strings.Contains(strings.ToLower(serviceName), filter) {
			continue
		}

		serviceNames = append(serviceNames, serviceName)
	}

	sort.Strings(serviceNames)

	return serviceNames
}

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

