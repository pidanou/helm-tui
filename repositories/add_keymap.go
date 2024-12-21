package repositories

import "github.com/charmbracelet/bubbles/key"

var addKeys = keyMap{
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}
