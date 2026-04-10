package keymaps

import "github.com/charmbracelet/bubbles/key"

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type ReleaseKeyMap struct {
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
func (k ReleaseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Install, k.Delete, k.Upgrade, k.Select, k.Refresh, k.Rollback, k.ChangeTab, k.Cancel, k.Back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k ReleaseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var releasesKeys = ReleaseKeyMap{
	Install: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install new release")),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Delete release"),
	),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")),
	Select:  key.NewBinding(key.WithKeys("enter/space"), key.WithHelp("enter/space", "Details")),
	Upgrade: key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
}

var historyKeys = ReleaseKeyMap{
	Install:  key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install new release")),
	Rollback: key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "Rollback to revision")),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Delete release"),
	),
	Upgrade:   key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
	ChangeTab: key.NewBinding(key.WithKeys("h", "l", "right", "left"), key.WithHelp("hl/←→", "Navigate tabs")),
	Back:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Back")),
}

var readOnlyKeys = ReleaseKeyMap{
	Install: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install new release")),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "Delete release"),
	),
	Upgrade:   key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade release")),
	ChangeTab: key.NewBinding(key.WithKeys("h", "l", "right", "left"), key.WithHelp("hl/←→", "Navigate tabs")),
	Back:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Back")),
}

func GenerateReleaseKeys() []ReleaseKeyMap {
	return []ReleaseKeyMap{releasesKeys, historyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys, readOnlyKeys}
}

var ReleaseInstallKeys = ReleaseKeyMap{
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}

var ReleaseUpgradeKeys = ReleaseKeyMap{
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}
