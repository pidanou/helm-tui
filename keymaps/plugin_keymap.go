package keymaps

import "github.com/charmbracelet/bubbles/key"

type PluginKeyMap struct {
	Install   key.Binding
	Update    key.Binding
	Uninstall key.Binding
	Cancel    key.Binding
	Refresh   key.Binding
}

var PluginOverviewKeys = PluginKeyMap{
	Uninstall: key.NewBinding(key.WithKeys(DefaultConfig.Plugin.Uninstall...), key.WithHelp(formatHelp(DefaultConfig.Plugin.Update), "Uninstall")),
	Install:   key.NewBinding(key.WithKeys(DefaultConfig.Plugin.Install...), key.WithHelp(formatHelp(DefaultConfig.Plugin.Install), "Install")),
	Update:    key.NewBinding(key.WithKeys(DefaultConfig.Plugin.Update...), key.WithHelp(formatHelp(DefaultConfig.Plugin.Update), "Update")),
	Cancel:    key.NewBinding(key.WithKeys(DefaultConfig.Plugin.Cancel...), key.WithHelp(formatHelp(DefaultConfig.Plugin.Cancel), "Cancel")),
	Refresh:   key.NewBinding(key.WithKeys(DefaultConfig.Plugin.Refresh...), key.WithHelp(formatHelp(DefaultConfig.Plugin.Refresh), "Refresh")),
}

func (k PluginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Update, k.Install, k.Uninstall, k.Refresh, k.Cancel}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k PluginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
