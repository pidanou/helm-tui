package plugins

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/types"
)

func (m PluginsModel) list() tea.Msg {
	var stdout bytes.Buffer
	var rows = []table.Row{}

	// Create the command
	cmd := exec.Command("helm", "plugin", "ls")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.ListReleasesMsg{Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		fields := strings.Fields(line)
		name := fields[0]
		version := fields[1]
		description := strings.Join(fields[2:], " ")
		row := []string{name, version, description}
		rows = append(rows, row)
	}
	return types.PluginsListMsg{Content: rows}
}

func (m PluginsModel) install() tea.Msg {
	pluginName := m.installPluginInput.Value()
	if pluginName == "" {
		return types.PluginInstallMsg{Err: errors.New("No plugin")}
	}
	cmd := exec.Command("helm", "plugin", "install", strings.TrimSpace(pluginName))
	err := cmd.Run()
	if err != nil {
		return types.PluginInstallMsg{Err: errors.New("Cannot install plugin")}
	}
	return types.PluginInstallMsg{Err: nil}
}

func (m PluginsModel) update() tea.Msg {
	if m.pluginsTable.SelectedRow() == nil {
		return types.PluginUpdateMsg{Err: errors.New("No plugin selected")}
	}
	pluginName := m.pluginsTable.SelectedRow()[0]
	cmd := exec.Command("helm", "plugin", "update", pluginName)
	err := cmd.Run()

	if err != nil {
		return types.PluginUpdateMsg{Err: errors.New("Cannot update plugin")}
	}
	return types.PluginUpdateMsg{Err: nil}
}

func (m PluginsModel) uninstall() tea.Msg {
	if m.pluginsTable.SelectedRow() == nil {
		return types.PluginUninstallMsg{Err: errors.New("No plugin selected")}
	}
	pluginName := m.pluginsTable.SelectedRow()[0]
	cmd := exec.Command("helm", "plugin", "uninstall", strings.TrimSpace(pluginName))
	err := cmd.Run()
	if err != nil {
		return types.PluginUninstallMsg{Err: errors.New("Cannot update plugin")}
	}
	return types.PluginUninstallMsg{Err: nil}
}
