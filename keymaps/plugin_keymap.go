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
	Uninstall: key.NewBinding(key.WithKeys("U"), key.WithHelp("U", "Uninstall")),
	Install:   key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install")),
	Update:    key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Update")),
	Cancel:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
	Refresh:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")),
}

func (k PluginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Update, k.Install, k.Uninstall, k.Refresh, k.Cancel}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k PluginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
