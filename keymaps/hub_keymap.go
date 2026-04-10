package keymaps

import (
	"github.com/charmbracelet/bubbles/key"
)

type HubKeyMap struct {
	StartSearch  key.Binding
	StartAddRepo key.Binding
	SendInput    key.Binding
	ShowValues   key.Binding
	Cancel       key.Binding
}

func (k HubKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.StartSearch, k.SendInput, k.StartAddRepo, k.ShowValues, k.Cancel}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k HubKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var HubDefaultKeysHelp = HubKeyMap{
	StartSearch: key.NewBinding(key.WithKeys(DefaultConfig.Hub.StartSearch...), key.WithHelp(formatHelp(DefaultConfig.Hub.StartSearch), "Search")),
}

var HubTableFocusedKeysHelp = HubKeyMap{
	ShowValues:  key.NewBinding(key.WithKeys(DefaultConfig.Hub.ShowValues...), key.WithHelp(formatHelp(DefaultConfig.Hub.ShowValues), "Show default values")),
	StartSearch: key.NewBinding(key.WithKeys(DefaultConfig.Hub.StartSearch...), key.WithHelp(formatHelp(DefaultConfig.Hub.SendInput), "Search")),
	SendInput:   key.NewBinding(key.WithKeys(DefaultConfig.Hub.SendInput...), key.WithHelp(formatHelp(DefaultConfig.Hub.SendInput), "Add repo")),
}

var HubSearchFocusedKeyHelp = HubKeyMap{
	SendInput: key.NewBinding(key.WithKeys(DefaultConfig.Hub.SendInput...), key.WithHelp(formatHelp(DefaultConfig.Hub.SendInput), "Search")),
}

var HubAddRepoFocusedKeyHelp = HubKeyMap{
	SendInput: key.NewBinding(key.WithKeys(DefaultConfig.Hub.SendInput...), key.WithHelp(formatHelp(DefaultConfig.Hub.SendInput), "Add Repo")),
}

var HubPackageDefaultValueKeyHelp = HubKeyMap{
	Cancel: key.NewBinding(key.WithKeys(DefaultConfig.Hub.Cancel...), key.WithHelp(formatHelp(DefaultConfig.Hub.Cancel), "Cancel")),
}
