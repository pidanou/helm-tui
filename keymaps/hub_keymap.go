package keymaps

import "github.com/charmbracelet/bubbles/key"

type HubKeyMap struct {
	AddRepo key.Binding
	Search  key.Binding
	Show    key.Binding
	Cancel  key.Binding
}

func (k HubKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.AddRepo, k.Show, k.Search, k.Cancel}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k HubKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var HubDefaultKeysHelp = HubKeyMap{
	Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Search")),
	Show:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Focus table")),
}

var HubTableKeysHelp = HubKeyMap{
	Show:    key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "Show default values")),
	Search:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Search")),
	AddRepo: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "Add repo")),
}

var HubSearchKeyHelp = HubKeyMap{
	Search: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Search")),
}

var HubAddRepoKeyHelp = HubKeyMap{
	Search: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Search")),
}

var HubDefaultValuesKeyHelp = HubKeyMap{
	Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Search")),
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}
