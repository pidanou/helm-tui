package repositories

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
)

func (m Model) View() string {
	helpView := m.help.View(m.keys[m.selectedView])
	repoView := m.renderTable(m.tables[listView], " Repositories ", m.selectedView == listView)
	packagesView := m.renderTable(m.tables[packagesView], " Packages ", m.selectedView == packagesView)
	versionsView := m.renderTable(m.tables[versionsView], " Versions ", m.selectedView == versionsView)
	view := lipgloss.JoinHorizontal(lipgloss.Top, repoView, packagesView, versionsView)
	if m.installing {
		return m.installModel.View()
	}
	if m.adding {
		return m.addModel.View()
	}
	if m.showDefaultValue {
		return m.renderDefaultValueView()
	}
	helperStyle := m.help.Styles.ShortSeparator
	return view + "\n" + helpView + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
}

func (m Model) renderTable(table table.Model, title string, active bool) string {
	var topBorder string
	table.SetHeight(m.height - 3)
	tableView := table.View()
	var baseStyle lipgloss.Style
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	topBorder = styles.GenerateTopBorderWithTitle(title, table.Width(), styles.Border, styles.InactiveStyle)
	if active {
		topBorder = styles.GenerateTopBorderWithTitle(title, table.Width(), styles.Border, styles.ActiveStyle.Foreground(styles.HighlightColor))
		baseStyle = styles.ActiveStyle.Border(styles.Border, false, true, true)
	}
	tableView = baseStyle.Render(tableView)
	return lipgloss.JoinVertical(lipgloss.Top, topBorder, tableView)
}

func (m Model) renderDefaultValueView() string {
	m.defaultValueVP.Height = m.height - 2 - 1
	defaultValueTopBorder := styles.GenerateTopBorderWithTitle(" Default Values ", m.defaultValueVP.Width, styles.Border, styles.InactiveStyle)
	baseStyle := styles.InactiveStyle.Border(styles.Border, false, true, true)
	helperStyle := m.help.Styles.ShortSeparator
	helpView := helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	return lipgloss.JoinVertical(lipgloss.Top, defaultValueTopBorder, baseStyle.Render(m.defaultValueVP.View()), m.help.View(defaultValuesKeyHelp)+helpView)
}
