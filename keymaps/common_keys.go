package keymaps

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type CommonKeyMap struct {
	MenuNext   key.Binding
	Quit       key.Binding
	Suggestion SuggestionKeyMap
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k CommonKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.MenuNext, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k CommonKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var menuKeys = append(DefaultConfig.Common.MenuPrev, DefaultConfig.Common.MenuNext...)
var CommonKeysHelper = CommonKeyMap{
	MenuNext:   key.NewBinding(key.WithKeys(menuKeys...), key.WithHelp(fmt.Sprintf("%s/%s", formatHelp(DefaultConfig.Common.MenuPrev), formatHelp(DefaultConfig.Common.MenuNext)), "Change panel")),
	Quit:       key.NewBinding(key.WithKeys(DefaultConfig.Common.Exit...), key.WithHelp(formatHelp(DefaultConfig.Common.Exit), "Quit")),
	Suggestion: SuggestionInputKeyMap,
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
