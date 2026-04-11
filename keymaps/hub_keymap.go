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

func HubDefaultKeysHelp() HubKeyMap {
	return HubKeyMap{
		StartSearch: key.NewBinding(key.WithKeys(DefaultConfig.Hub.StartSearch...), key.WithHelp(formatHelp(DefaultConfig.Hub.StartSearch), "Search")),
	}
}

func HubTableFocusedKeysHelp() HubKeyMap {
	return HubKeyMap{
		ShowValues:  key.NewBinding(key.WithKeys(DefaultConfig.Hub.ShowValues...), key.WithHelp(formatHelp(DefaultConfig.Hub.ShowValues), "Show default values")),
		StartSearch: key.NewBinding(key.WithKeys(DefaultConfig.Hub.StartSearch...), key.WithHelp(formatHelp(DefaultConfig.Hub.SendInput), "Search")),
		SendInput:   key.NewBinding(key.WithKeys(DefaultConfig.Hub.SendInput...), key.WithHelp(formatHelp(DefaultConfig.Hub.SendInput), "Add repo")),
	}
}

func HubSearchFocusedKeyHelp() HubKeyMap {
	return HubKeyMap{
		SendInput: key.NewBinding(key.WithKeys(DefaultConfig.Hub.SendInput...), key.WithHelp(formatHelp(DefaultConfig.Hub.SendInput), "Search")),
	}
}

func HubAddRepoFocusedKeyHelp() HubKeyMap {
	return HubKeyMap{
		SendInput: key.NewBinding(key.WithKeys(DefaultConfig.Hub.SendInput...), key.WithHelp(formatHelp(DefaultConfig.Hub.SendInput), "Add Repo")),
	}
}

func HubPackageDefaultValueKeyHelp() HubKeyMap {
	return HubKeyMap{
		Cancel: key.NewBinding(key.WithKeys(DefaultConfig.Hub.Cancel...), key.WithHelp(formatHelp(DefaultConfig.Hub.Cancel), "Cancel")),
	}
}
