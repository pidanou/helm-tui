package plugins

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/components"
	"github.com/pidanou/helmtui/types"
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
	keys               keyMap
	width              int
	height             int
}

func InitModel() PluginsModel {
	table := components.GenerateTable()
	input := textinput.New()
	input.Placeholder = "Enter plugin path/url"
	return PluginsModel{pluginsTable: table, help: help.New(), keys: overviewKeys, installPluginInput: textinput.New()}
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
		return m, m.list
	case types.PluginUninstallMsg:
		return m, m.list
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Install):
			cmds = append(cmds, m.installPluginInput.Focus())
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.Uninstall):
			if !m.installPluginInput.Focused() {
				return m, m.uninstall
			}
		case key.Matches(msg, m.keys.Update):
			if !m.installPluginInput.Focused() {
				return m, m.update
			}
		case key.Matches(msg, m.keys.Cancel):
			m.installPluginInput.Blur()
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.Refresh):
			return m, m.list
		case msg.String() == "enter":
			if m.installPluginInput.Focused() {
				return m, m.install
			}
		}
	}
	m.installPluginInput, cmd = m.installPluginInput.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
