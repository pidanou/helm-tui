package hub

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
)

func (m HubModel) View() string {
	header := styles.InactiveStyle.Border(styles.Border).Render(m.searchBar.View())
	remainingHeight := m.height - lipgloss.Height(header) - 2 - 1 //  searchbar padding + releaseTable padding + helper
	if m.repoAddInput.Focused() {
		remainingHeight -= 3
	}
	if m.view == defaultValueView {
		m.defaultValueVP.Height = m.height - 2 - 1
		return m.renderDefaultValueView()
	}
	m.resultTable.SetHeight(remainingHeight)
	if m.searchBar.Focused() {
		header = styles.ActiveStyle.Border(styles.Border).Render(m.searchBar.View())
	}
	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(defaultKeysHelp)
	if m.searchBar.Focused() {
		helpView = m.help.View(searchKeyHelp)
	}
	if m.resultTable.Focused() {
		helpView = m.help.View(tableKeysHelp)
	}
	if m.repoAddInput.Focused() {
		helpView = m.help.View(addRepoKeyHelp)
	}
	helpView += helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	style := styles.ActiveStyle.Border(styles.Border)
	if m.repoAddInput.Focused() {
		return header + "\n" + m.renderSearchTableView() + "\n" + style.Render(m.repoAddInput.View()) + "\n" + helpView
	}
	return header + "\n" + m.renderSearchTableView() + "\n" + helpView
}

func (m HubModel) renderSearchTableView() string {
	var releasesTopBorder string
	tableView := m.resultTable.View()
	var baseStyle lipgloss.Style
	releasesTopBorder = styles.GenerateTopBorderWithTitle(" Results ", m.resultTable.Width(), styles.Border, styles.InactiveStyle)
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	if m.resultTable.Focused() {
		releasesTopBorder = styles.GenerateTopBorderWithTitle(" Results ", m.resultTable.Width(), styles.Border, styles.ActiveStyle.Foreground(styles.HighlightColor))
		baseStyle = styles.ActiveStyle.Border(styles.Border, false, true, true)
	}
	tableView = baseStyle.Render(tableView)
	return lipgloss.JoinVertical(lipgloss.Top, releasesTopBorder, tableView)
}

func (m HubModel) renderDefaultValueView() string {
	defaultValueTopBorder := styles.GenerateTopBorderWithTitle(" Default Values ", m.defaultValueVP.Width, styles.Border, styles.InactiveStyle)
	baseStyle := styles.InactiveStyle.Border(styles.Border, false, true, true)
	helperStyle := m.help.Styles.ShortSeparator
	helpView := helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	return lipgloss.JoinVertical(lipgloss.Top, defaultValueTopBorder, baseStyle.Render(m.defaultValueVP.View()), m.help.View(defaultValuesKeyHelp)+helpView)
}
