package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

const (
	defaultTableHeight = 20
	minTableHeight     = 5
	verticalChromeRows = 7
)

type Model struct {
	systemdManager    *systemd.SystemdManager
	err               error
	interval          time.Duration
	selectedServiceID string
	actionInFlight    bool
	pendingAction     systemd.ServiceAction
	pendingServiceID  string
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
		tableHeight:    defaultTableHeight,
		monitoring:     true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		monitorStateCmd(m.systemdManager),
		tick(m.interval),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resizeTable(msg.Height)
		m.normalizeSelection()
		return m, nil

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
		m.actionInFlight = false
		m.pendingAction = ""
		m.pendingServiceID = ""
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
				if m.actionInFlight {
					return m, nil
				}

				if m.selectedServiceID == "" {
					return m, nil
				}

				m.actionInFlight = true
				m.pendingAction = action
				m.pendingServiceID = m.selectedServiceID

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
