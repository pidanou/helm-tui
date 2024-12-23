package helpers

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	MenuNext key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MenuNext, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var CommonKeys = keyMap{
	MenuNext: key.NewBinding(key.WithKeys("[", "]"), key.WithHelp("[/]", "Change panel")),
	Quit:     key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "Quit")),
}

type SuggestionKeyMap struct {
	AcceptSuggestion key.Binding
	NextSuggestion   key.Binding
	PrevSuggestion   key.Binding
}

var SuggestionInputKeyMap = SuggestionKeyMap{
	AcceptSuggestion: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "Accept suggestion")),
	NextSuggestion:   key.NewBinding(key.WithKeys("down", "ctrl+n"), key.WithHelp("down/ctrl+n", "Next suggestion")),
	PrevSuggestion:   key.NewBinding(key.WithKeys("up", "ctrl+p"), key.WithHelp("up/ctrl+p", "Previous suggestion")),
}

func (k SuggestionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.AcceptSuggestion, k.NextSuggestion, k.PrevSuggestion}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k SuggestionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
