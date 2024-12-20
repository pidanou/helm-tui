package releases

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Install   key.Binding
	Delete    key.Binding
	Rollback  key.Binding
	Refresh   key.Binding
	Select    key.Binding
	ChangeTab key.Binding
	Back      key.Binding
	Upgrade   key.Binding
	Cancel    key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Install, k.Delete, k.Upgrade, k.Select, k.Refresh, k.Rollback, k.ChangeTab, k.Cancel, k.Back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var releasesKeys = keyMap{
	Install: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install new release")),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete"),
	),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")),
	Select:  key.NewBinding(key.WithKeys("enter/space"), key.WithHelp("enter/space", "Details")),
	Upgrade: key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
}

var historyKeys = keyMap{
	Install:  key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install new release")),
	Rollback: key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "Rollback to revision")),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete"),
	),
	Upgrade:   key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
	ChangeTab: key.NewBinding(key.WithKeys("tab", "shift+tab"), key.WithHelp("tab/shift+tab", "Navigate tabs")),
	Back:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Back")),
}

var readOnlyKeys = keyMap{
	Install: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install new release")),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete release"),
	),
	Upgrade:   key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
	ChangeTab: key.NewBinding(key.WithKeys("tab", "shift+tab"), key.WithHelp("tab/shift+tab", "Navigate tabs")),
	Back:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Back")),
}
var installKeys = keyMap{
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}

func generateKeys() []keyMap {
	return []keyMap{releasesKeys, historyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys}
}
