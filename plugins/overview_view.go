package plugins

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helm-tui/components"
	"github.com/pidanou/helm-tui/keymaps"
	"github.com/pidanou/helm-tui/styles"
)

func (m PluginsModel) View() string {
	var remainingHeight = m.height
	if m.installPluginInput.Focused() {
		remainingHeight -= 3
	}
	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(keymaps.PluginOverviewKeys) + helperStyle.Render(" • ") + m.help.View(keymaps.CommonKeysHelper)
	view := components.RenderTable(m.pluginsTable, remainingHeight-3, m.width-2)
	m.installPluginInput.Width = m.width - 5
	if m.installPluginInput.Focused() {
		view += "\n" + styles.ActiveStyle.Border(styles.Border).Render(m.installPluginInput.View())
	}
	view = lipgloss.JoinVertical(lipgloss.Left, view, helpView)
	return view
}
