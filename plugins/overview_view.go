package plugins

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/components"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
)

func (m PluginsModel) View() string {
	var remainingHeight = m.height
	if m.installPluginInput.Focused() {
		remainingHeight -= 3
	}
	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(m.keys) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	view := components.RenderTable(m.pluginsTable, remainingHeight-3, m.width-2)
	m.installPluginInput.Width = m.width - 5
	if m.installPluginInput.Focused() {
		view += "\n" + styles.ActiveStyle.Border(styles.Border).Render(m.installPluginInput.View())
	}
	view = lipgloss.JoinVertical(lipgloss.Left, view, helpView)
	return view
}
