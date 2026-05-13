package tui

import (
	"sort"
	"strings"
)

func (m *Model) resizeTable(windowHeight int) {
	tableHeight := windowHeight - verticalChromeRows
	if tableHeight < minTableHeight {
		tableHeight = minTableHeight
	}

	m.tableHeight = tableHeight
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
