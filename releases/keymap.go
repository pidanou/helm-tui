package releases

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Delete    key.Binding
	Rollback  key.Binding
	Refresh   key.Binding
	Select    key.Binding
	ChangeTab key.Binding
	Back      key.Binding
	Upgrade   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Delete, k.Upgrade, k.Select, k.Refresh, k.Rollback, k.ChangeTab, k.Back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var releasesKeys = keyMap{
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete"),
	),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")),
	Select:  key.NewBinding(key.WithKeys("enter/space"), key.WithHelp("enter/space", "Details")),
	Upgrade: key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
}

var historyKeys = keyMap{
	Rollback: key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "Rollback to revision")),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete"),
	),
	Upgrade:   key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
	ChangeTab: key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←,h/→,l", "Navigate tabs")),
	Back:      key.NewBinding(key.WithKeys("esc/backspace"), key.WithHelp("esc/backspace", "Back")),
}

var readOnlyKeys = keyMap{
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete release"),
	),
	Upgrade:   key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
	ChangeTab: key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←,h/→,l", "Navigate tabs")),
	Back:      key.NewBinding(key.WithKeys("esc/backspace"), key.WithHelp("esc/backspace", "Back")),
}

func generateKeys() []keyMap {
	return []keyMap{releasesKeys, historyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys}
}
