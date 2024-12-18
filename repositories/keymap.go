package repositories

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Delete  key.Binding
	Refresh key.Binding
	Move    key.Binding
	Update  key.Binding
	Install key.Binding
	Select  key.Binding
	Cancel  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Delete, k.Update, k.Move, k.Select, k.Refresh, k.Install, k.Cancel}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var repoListKeys = keyMap{
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete repo"),
	),
	Move:    key.NewBinding(key.WithKeys("h", "	j", "k", "l", "left", "right", "up", "down"), key.WithHelp("hjkl/↑/↓/←/→", "Move")),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Refresh")),
	Select:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Select")),
	Update:  key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Update repo")),
	Install: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install version")),
}

var chartsListKeys = keyMap{
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete repo"),
	),
	Move:    key.NewBinding(key.WithKeys("h", "	j", "k", "l", "left", "right", "up", "down"), key.WithHelp("hjkl/↑/↓/←/→", "Move")),
	Select:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Select")),
	Update:  key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Update repo")),
	Install: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install version")),
}

var versionsKeys = keyMap{
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete repo"),
	),
	Move:    key.NewBinding(key.WithKeys("h", "	j", "k", "l", "left", "right", "up", "down"), key.WithHelp("hjkl/↑/↓/←/→", "Move")),
	Update:  key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "Upgrade repo")),
	Install: key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "Install version")),
}

var installKeys = keyMap{
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}

func generateKeys() []keyMap {
	return []keyMap{repoListKeys, chartsListKeys, versionsKeys}
}
