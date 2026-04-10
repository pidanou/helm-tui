package plugins

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/components"
	"github.com/pidanou/helm-tui/keymaps"
	"github.com/pidanou/helm-tui/types"
)

var pluginsCols = []components.ColumnDefinition{
	{Title: "Name", FlexFactor: 1},
	{Title: "Version", FlexFactor: 1},
	{Title: "description", FlexFactor: 3},
}

type PluginsModel struct {
	pluginsTable       table.Model
	installPluginInput textinput.Model
	help               help.Model
	width              int
	height             int
}

func InitModel() PluginsModel {
	table := components.GenerateTable()
	input := textinput.New()
	input.Placeholder = "Enter plugin path/url"
	return PluginsModel{
		pluginsTable:       table,
		help:               help.New(),
		installPluginInput: input}
}

func (m PluginsModel) Init() tea.Cmd {
	return m.list
}

func (m PluginsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		components.SetTable(&m.pluginsTable, pluginsCols, m.width)
	case types.PluginsListMsg:
		m.pluginsTable.SetRows(msg.Content)
	case types.PluginInstallMsg:
		m.installPluginInput.Blur()
		m.installPluginInput.SetValue("")
		return m, m.list
	case types.PluginUninstallMsg:
		return m, m.list
	case tea.KeyMsg:
		switch {
		case keymaps.Contains(keymaps.DefaultConfig.Plugin.Install, msg.String()):
			cmds = append(cmds, m.installPluginInput.Focus())
			return m, tea.Batch(cmds...)
		case keymaps.Contains(keymaps.DefaultConfig.Plugin.Uninstall, msg.String()):
			if !m.installPluginInput.Focused() {
				return m, m.uninstall
			}
		case keymaps.Contains(keymaps.DefaultConfig.Plugin.Update, msg.String()):
			if !m.installPluginInput.Focused() {
				return m, m.update
			}
		case keymaps.Contains(keymaps.DefaultConfig.Plugin.Cancel, msg.String()):
			m.installPluginInput.Blur()
			return m, tea.Batch(cmds...)
		case keymaps.Contains(keymaps.DefaultConfig.Plugin.Refresh, msg.String()):
			return m, m.list
		case keymaps.Contains(keymaps.DefaultConfig.Plugin.ConfirmInstall, msg.String()):
			if m.installPluginInput.Focused() {
				return m, m.install
			}
		}
	}
	m.installPluginInput, cmd = m.installPluginInput.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
