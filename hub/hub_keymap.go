package hub

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	AddRepo key.Binding
	Search  key.Binding
	Show    key.Binding
	Cancel  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.AddRepo, k.Show, k.Search, k.Cancel}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var defaultKeysHelp = keyMap{
	Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Search")),
	Show:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Focus table")),
}

var tableKeysHelp = keyMap{
	Show:    key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "Show default values")),
	Search:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Search")),
	AddRepo: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "Add repo")),
}

var searchKeyHelp = keyMap{
	Search: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Search")),
}

var addRepoKeyHelp = keyMap{
	Search: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Search")),
}

var defaultValuesKeyHelp = keyMap{
	Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Search")),
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}
